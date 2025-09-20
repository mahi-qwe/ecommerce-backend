package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
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

	r.Use(middlewares.CORSMiddleware())

	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	routes.ProductRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-commerce API running!", //sample endpoint
		})
	})

	r.Run(":8080")
}
