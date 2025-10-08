package main

import (
	pb "audio-notes/go-source"
	"audio-notes/models"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func SetupRouter(uploader *manager.Uploader, s3Client *s3.Client) *gin.Engine {
	r := gin.Default()

	// create slice of notes
	var notes []models.Note

	// get endpoint
	r.GET("/greet", func(c *gin.Context) {
		// c.String(http.StatusOK, "Welcome to your first Go API!")
		c.JSON(http.StatusOK, gin.H{
			"message": "Greetings!",
		})
	})
	// post endpoint
	r.POST("/upload", func(c *gin.Context) {
		// Get file from form
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Open the file
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer f.Close()

		// Upload to S3
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("s3-golang-uploaded-pdfs"),
			Key:    aws.String(file.Filename),
			Body:   f,
		})
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
			return
		}

		fmt.Println("Successfully uploaded to:", result.Location)

		// Save note info
		newNote := models.Note{
			ID:         len(notes) + 1,
			FilePath:   result.Location, // use S3 URL
			Title:      "My New Note!",
			UploadTime: time.Now(),
		}
		notes = append(notes, newNote)

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully uploaded!",
			"note":    newNote,
		})
	})
	r.GET("/bucket-files", func(c *gin.Context) {
		paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String("s3-golang-uploaded-pdfs"),
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				fmt.Print("Error! ", err)
			}
			// prints out bucket name
			log.Printf("Bucket Name: ", *page.Name)
			for index, item := range page.Contents {
				// item.Key is a pointer to the string, so need to deref
				fmt.Println(index, *item.Key)
			}
		}
	})
	r.GET("/all_notes", func(c *gin.Context) {
		for i:=0; i < len(notes); i++ {
			note := notes[i]
			fmt.Print(note)
		}
		c.JSON(http.StatusOK, gin.H{
			"notes": notes,
		})
	})
	r.GET("/narrate-notes", func(c *gin.Context) {
		/*
			- Make a request with URL to eleven labs
			- Add headers to API call (xi-api-key), voice_id
		*/
		headers := http.Header{}
		headers.Add("xi-api-key", os.Getenv("ELEVEN_LABS_KEY"))
		// pass in voice id to the api; used GET endpoint from docs
		url := fmt.Sprintf("wss://api.elevenlabs.io/v1/text-to-speech/%s/stream-input", "29vD33N1CtxCmqQRPOHJ")
		// Dial returns connection, response (optional, we don't need as of now - can use to inspect headers, other handshake elements), and an error
		connection, _, err := websocket.DefaultDialer.Dial(url, headers)
		if err != nil {
			fmt.Println("Error establishing connection to ElevenLabs WebSocket", err)
		}

		// call the Python server with a URL to the S3 file
		// get response
		// feeed into eleven labs api
		// test - s3://s3-golang-uploaded-pdfs/Biology Notes.pdf

		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Error establishing connection with client! ", err)
		}
		defer conn.Close()
		client := pb.NewTextExtractionServiceClient(conn)
		req, err := client.ExtractText(context.Background(), &pb.ExtractTextRequest{
			S3Url: "s3://s3-golang-uploaded-pdfs/Biology Notes.pdf",
		})
		if err != nil {
			fmt.Println("Error making request: ", err)
		}
		fmt.Println(req.ExtractedText)
		

		testJson := map[string]interface{} {
			"type" : "sendText",
			"text" : req.ExtractedText,
			"flush" : true,
		}
		// reassigning so we can omit the :
		err = connection.WriteJSON(testJson)
		if err != nil {
			fmt.Println("Error writing json to server: ", err)
		}

		for {
			// omitted (_) is messageType an integer
			// tells us what type of message was received (ping, pong, binary, text, close)
			_, bytes, err := connection.ReadMessage()
			if err != nil {
				fmt.Println("Connection esablished, but error reading message: ", err)
				break
			}
			// we declare a map to hold the decoded string coming back from our api request
			var audioResponseMapping map[string]interface{}
			// unmarshal: takes array of bytes and destination (pointer) to where it needs to be stored 
			err = json.Unmarshal(bytes, &audioResponseMapping)
			if err != nil {
				fmt.Println("Error unmarshalling: ", err)
			}
			// base64 decoded string representing our audio
			// when we unmarshall, every value in the mapping becomes of type interface{}, so we need to assert it as a string
			audioDecodedString, ok := audioResponseMapping["audio"].(string)
			if !ok {
				fmt.Println("Decoded Output is not a string: ", ok)
			}
			decodedByteResponse, err := base64.StdEncoding.DecodeString(audioDecodedString)
			if err != nil {
				fmt.Println("Error when trying to decode the base64 output")
			}
			err = os.WriteFile("testing.mp3", decodedByteResponse, 0100644)
			if err != nil {
				fmt.Println("Error saving file!", err)
			}
		}
		

		// defer will inherently close the connection once our callback function is complete, so we don't have to remember about it
		defer connection.Close()
	})
	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/*
	- config is a Go package that provides functions and types to load and manage AWS Configuration
		- helps program let you know what credentials to use, region, etc
	- LoadDefaultConfig() is a function from the AWS Go SDK
		- loads the AWS Configuration, including credentials, default region, etc
	- context.TODO()
		- a lot of the AWS SDK functions require a context.Context
		- placeholder context since we don't have specific context just yet
	*/
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	/*
	- client is a Go object that knows credentials, region, and AWS endpoints
	- all requests trying to hit S3 go through this client object
	*/
	client := s3.NewFromConfig(cfg)
	/*
	- S3 objects can be uploaded in parts; the uploader is a convenient wrapper that does that automatically
	*/
	uploader := manager.NewUploader(client)


	r := SetupRouter(uploader, client)
	r.Run(":3000")
}
