package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// ✅ POST /wishlist - Add product to wishlist
func AddToWishlist(c *gin.Context) {
	userID := getUserID(c)

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check product exists
	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Prevent duplicates
	var existing models.WishlistItem
	if err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in wishlist"})
		return
	}

	// Create wishlist item
	item := models.WishlistItem{
		UserID:    userID,
		ProductID: input.ProductID,
	}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	// Preload product
	config.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added to wishlist",
		"item":    item,
	})
}

// ✅ GET /wishlist - View wishlist items
func GetWishlist(c *gin.Context) {
	userID := getUserID(c)

	var items []models.WishlistItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": items})
}

// ✅ DELETE /wishlist/:product_id - Remove item from wishlist
func RemoveFromWishlist(c *gin.Context) {
	userID := getUserID(c)
	productID := c.Param("product_id")

	result := config.DB.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.WishlistItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist"})
}
