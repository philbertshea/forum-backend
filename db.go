package main

import (
	"forum-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Thread{}, &models.Comment{})
	log.Println("Database Migration done")
}

// User CRUD Code
func createUser(c *gin.Context) {
	var user models.User // Create a new user
	// Check for errors
	err := c.ShouldBindJSON(&user) // Bind to JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Check if user with similar username exists
	var existingUser models.User
	err = DB.Where("username = ?", user.Username).First(&existingUser).Error
	// If a user does NOT exist, there should BE an error. (No First Row - record not found)
	// If a user exists, there should be NO error.
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A user already exists."})
	}

	err = DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func getUsers(c *gin.Context) {
	var allUsers []models.User

	// Query all users from the database
	err := DB.Table("users").Find(&allUsers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of users
	c.JSON(http.StatusOK, allUsers)

}

func loginUser(c *gin.Context) {
	var requestBody struct {
		Username string `json:"username"`
	}

	err := c.ShouldBindJSON(&requestBody) // Bind to JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User
	err = DB.Where("username = ?", requestBody.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist."})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Successful"})
}

/*
func updateUser(c *gin.Context) {
	var existingUser models.User
	// Check for errors
	err := c.ShouldBindJSON(&existingUser) // Bind to JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Check if user with similar username exists
	err = DB.Where("username = ?", existingUser.Username).First(&existingUser).Error
	// If a user exists, there should be NO error.
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "No user with that username exists."})
	}
	existingUser

	c.JSON(http.StatusOK, user)
}
*/

func createThread(c *gin.Context) {
	var thread models.Thread // Create a new user
	// Check for errors
	err := c.ShouldBindJSON(&thread) // Bind to JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = DB.Create(&thread).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, thread)
}

func getThreads(c *gin.Context) {
	var threads []models.Thread // Create a new user

	err := DB.Table("threads").Find(&threads).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, threads)
}

func getIdFromUsername(c *gin.Context) {

	var requestBody struct {
		Username string `json:"username"`
	}

	err := c.ShouldBindJSON(&requestBody) // Bind to JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User

	err = DB.Where("username = ?", requestBody.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of users
	c.JSON(http.StatusOK, user.ID)

}

func getThreadFromId(c *gin.Context) {
	// Get id from user parameter
	id := c.Param("id")
	var thread models.Thread

	err := DB.Preload("Comments").First(&thread, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, thread)
}
