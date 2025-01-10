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
	var users []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	err := DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, users)

}
