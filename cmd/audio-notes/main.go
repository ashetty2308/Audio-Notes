package main

import(
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/greet", func(c *gin.Context) {
		// c.String(http.StatusOK, "Welcome to your first Go API!")
		c.JSON(http.StatusOK, gin.H{
			"message" : "Greetings!",
		})
	})	
	return r
}

func main() {
	r := SetupRouter()
	r.Run(":3000")
}