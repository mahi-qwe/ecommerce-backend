package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mahi-qwe/ecommerce-backend/config" // import from root module
	"github.com/mahi-qwe/ecommerce-backend/models"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect DB
	config.ConnectDatabase()

	// Use = instead of := to reuse err
	err = config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Init Gin
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-commerce API running 🚀",
		})
	})

	r.Run(":8080")
}
