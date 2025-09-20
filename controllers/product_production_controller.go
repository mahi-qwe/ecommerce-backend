package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// StartProductionHandler handles POST /admin/products/:id/production, also fetches the associated product
func StartProductionHandler(c *gin.Context) {
	productID := c.Param("id")

	// Check if product exists
	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if production already exists for this product
	var existingProduction models.ProductProduction
	if err := config.DB.Where("product_id = ? AND deleted_at IS NULL", productID).First(&existingProduction).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Production already started for this product"})
		return
	}

	// Create new production record
	production := models.ProductProduction{
		ProductID: product.ID,
		Status:    "pending",
		StartedAt: time.Now(),
	}

	if err := config.DB.Create(&production).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start production"})
		return
	}

	// Preload the Product relation
	if err := config.DB.Preload("Product").First(&production, production.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch production with product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Production started successfully",
		"production": production,
	})
}
