package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"gorm.io/gorm"
)

type AddCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func AddToCart(c *gin.Context) {
	// Get userID safely
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDValue.(type) {
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	// Bind JSON
	var input AddCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists and is not soft deleted
	var product models.Product
	if err := config.DB.Where("id = ? AND deleted_at IS NULL", input.ProductID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or deleted"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if cart item already exists
	var cartItem models.CartItem
	err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cartItem).Error

	if err == nil {
		// Already in cart → check stock
		newQuantity := cartItem.Quantity + input.Quantity
		if newQuantity > product.StockQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Exceeds available stock"})
			return
		}

		cartItem.Quantity = newQuantity
		if err := config.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
			return
		}

		// Preload associations
		if err := config.DB.Preload("User").Preload("Product").First(&cartItem, cartItem.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart updated", "cart_item": cartItem})
		return
	}

	// New cart item → check stock
	if input.Quantity > product.StockQuantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	newCartItem := models.CartItem{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := config.DB.Create(&newCartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	// Preload associations
	if err := config.DB.Preload("User").Preload("Product").First(&newCartItem, newCartItem.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Product added to cart",
		"cart_item": newCartItem,
	})
}
