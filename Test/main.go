package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, this is your API!"})
	})

	// Define a route for an example endpoint
	router.GET("/example", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is an example endpoint"})
	})

	// Run the server on port 8080
	router.Run(":8080")
}
