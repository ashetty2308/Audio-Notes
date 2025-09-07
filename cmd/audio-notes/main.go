package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"audio-notes/models"
	"time"
	"fmt"
	"encoding/json"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/greet", func(c *gin.Context) {
		// c.String(http.StatusOK, "Welcome to your first Go API!")
		c.JSON(http.StatusOK, gin.H{
			"message": "Greetings!",
		})
	})
	return r
}

func main() {
	// r := SetupRouter()
	// r.Run(":3000")

	note1 := models.Note{
		ID:	1,
		Title: "My First Note!",
		FilePath: "/uploads/lecture1.pdf",
		UploadTime: time.Now(),
	}
	note_as_json, err := json.Marshal(note1)
	if err != nil {
		fmt.Println("Failed to marshal note: ", err)
	}
	fmt.Println(string(note_as_json))
}
