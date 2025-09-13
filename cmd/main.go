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
	log.Println("✅ Users table migrated successfully")

	err = config.DB.AutoMigrate(&models.OTP{})
	if err != nil {
		log.Fatal("Migration failed for OTP table: ", err)
	}
	log.Println("✅ OTP table migrated successfully")

	err = config.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Migration failed for Product table: ", err)
	}
	log.Println("✅ Product table migrated successfully")

	err = config.DB.AutoMigrate(&models.ProductProduction{})
	if err != nil {
		log.Fatal("Migration failed for ProductProduction table: ", err)
	}
	log.Println("✅ ProductProduction table migrated successfully")

	err = config.DB.AutoMigrate(&models.CartItem{})
	if err != nil {
		log.Fatal("Migration failed for CartItem table: ", err)
	}
	log.Println("✅ CartItem table migrated successfully")

	// Init Gin
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-commerce API running 🚀",
		})
	})

	r.Run(":8080")
}
