package main

import (
	"net/http"

	_ "github.com/School-meal-lover/backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Grrrrr API
// @version 1.0
// @description The server for Grrrrr application.

func main() {
	router := gin.Default()

	// Set up the API routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	router.GET("/hello", Helloworld)

	// Set up Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	router.Run(":8080")
}

// @Summary Hello World
// @Description Returns a hello world message
// @Tags hello
// @Accept  json
// @Produce  json
// @Success 200 {string} string "hello world"
// @Router /hello [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
