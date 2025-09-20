package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// CreateProductHandler handles POST /admin/products
func CreateProductHandler(c *gin.Context) {
	var input models.Product

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": input,
	})
}

// GetProductsHandler handles GET /products (public)
func GetProductsHandler(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// GetProductByIDHandler handles GET /products/:id
func GetProductByIDHandler(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// UpdateProductHandler handles PUT /admin/products/:id
func UpdateProductHandler(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if product exists
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Bind JSON input
	var input struct {
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		Price         *float64 `json:"price"`
		StockQuantity *int     `json:"stock_quantity"`
		Category      *string  `json:"category"`
		ImageURL      *string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.StockQuantity != nil {
		updates["stock_quantity"] = *input.StockQuantity
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.ImageURL != nil {
		updates["image_url"] = *input.ImageURL
	}

	// Nothing to update
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	// Update timestamp
	updates["updated_at"] = time.Now()

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

// DeleteProductHandler handles DELETE /admin/products/:id
func DeleteProductHandler(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if product exists
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Soft delete product
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
