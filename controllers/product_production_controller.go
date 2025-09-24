package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// StartProductionHandler handles POST /admin/products/:id/production, also fetches/Preloads the associated product
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

// UpdateProductionStatusHandler handles PUT /admin/products/:id/production/status
func UpdateProductionStatusHandler(c *gin.Context) {
	productID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	// Validate status
	if req.Status != "pending" && req.Status != "in_progress" && req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Find production record for product
	var production models.ProductProduction
	if err := config.DB.Where("product_id = ?", productID).First(&production).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Production record not found"})
		return
	}

	// If status is completed, set CompletedAt
	if req.Status == "completed" {
		now := time.Now()
		production.CompletedAt = &now
	}

	production.Status = req.Status

	if err := config.DB.Save(&production).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update production status"})
		return
	}

	// Reload with Product preloaded
	if err := config.DB.Preload("Product").First(&production, production.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch production with product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Production status updated successfully",
		"production": production,
	})
}

// GetProductionDetailsHandler handles GET /admin/products/:id/production
func GetProductionDetailsHandler(c *gin.Context) {
	productID := c.Param("id")

	var production models.ProductProduction
	if err := config.DB.
		Preload("Product").
		Where("product_id = ?", productID).
		First(&production).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Production record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"production": production,
	})
}

// GetAllProductionsHandler handles GET /admin/products/production
func GetAllProductionsHandler(c *gin.Context) {
	var productions []models.ProductProduction

	if err := config.DB.
		Preload("Product").
		Find(&productions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch productions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"productions": productions,
	})
}
