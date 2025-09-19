package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"audio-notes/models"
	"time"
	"fmt"
	// "encoding/json"
)

func SetupRouter() *gin.Engine {
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
	r.POST("/post", func(c *gin.Context) {
		newNote := models.Note{
			ID: 3,
			FilePath: "temporary/file/path",
			Title: "My New Note!",
			UploadTime: time.Now(),
		}
		notes = append(notes, newNote)
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully created a new note!",
			"note": newNote,
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
	r := SetupRouter()
	r.Run(":3000")

	// note1 := models.Note{
	// 	ID:	1,
	// 	Title: "My First Note!",
	// 	FilePath: "/uploads/lecture1.pdf",
	// 	UploadTime: time.Now(),
	// }
	// note_as_json, err := json.Marshal(note1)
	// if err != nil {
	// 	fmt.Println("Failed to marshal note: ", err)
	// }
	// fmt.Println(string(note_as_json))
}
