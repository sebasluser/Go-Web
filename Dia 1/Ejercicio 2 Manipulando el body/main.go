package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Persona struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func main() {
	router := gin.Default()

	router.POST("/saludo", func(c *gin.Context) {
		var p Persona

		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		greeting := "Hola " + p.Nombre + " " + p.Apellido

		c.JSON(http.StatusOK, gin.H{"message": greeting})
	})

	router.Run(":8080")

}
