package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"github.com/mahi-qwe/ecommerce-backend/routes"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect DB
	config.ConnectDatabase()

	// Run migrations
	models.Migrate()

	r := gin.Default()

	routes.AuthRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-commerce API running!", //sample endpoint
		})
	})

	r.Run(":8080")
}
