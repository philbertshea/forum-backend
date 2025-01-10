package main

import (
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func sayHi(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Welcome to the forum backend.",
	})
}

func main() {
	// Initialise the Forum Database
	InitDB()

	router := gin.Default()

	// Enable CORS
	// Allow requests from the frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/", sayHi)

	// POST routes
	router.POST("/registerUser", createUser)
	router.POST("/getUsers", getUsers)

	// Run
	router.Run()
}
