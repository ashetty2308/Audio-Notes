package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"audio-notes/models"
	"time"
	"fmt"
	"context"
	"log"
	"github.com/joho/godotenv"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

func SetupRouter(uploader *manager.Uploader) *gin.Engine {
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
	r.GET("/all_notes", func(c *gin.Context) {
		for i:=0; i < len(notes); i++ {
			note := notes[i]
			fmt.Print(note)
		}
		c.JSON(http.StatusOK, gin.H{
			"notes": notes,
		})
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


	r := SetupRouter(uploader)
	r.Run(":3000")
}
