package main

import (
	"Practica/cmd/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	server := gin.Default()

	group := server.Group("/products")

	router := handler.NewProductRouter(group)

	router.ProductRoutes()

	server.Run(":8080")
}
