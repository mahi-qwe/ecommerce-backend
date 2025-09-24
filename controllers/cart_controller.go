package controllers

import (
	// "net/http" , i have used numbers directly instead of http.StatusOK etc
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	// "gorm.io/gorm"
)

type CartItemResponse struct {
	ID        uint           `json:"id"`
	Product   ProductSummary `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ProductSummary struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	ImageURL      string  `json:"image_url"`
}

func AddToCart(c *gin.Context) {
	userID := getUserID(c) // helper to safely get userID

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔹 Check product exists and not soft-deleted
	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 🔹 Check stock availability
	if input.Quantity > product.StockQuantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	// 🔹 Check if already in cart → STRICT separation (reject if exists)
	var existingItem models.CartItem
	if err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existingItem).Error; err == nil {
		// Already exists → reject, force client to use PUT
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in cart. Use PUT /cart/:id to update quantity."})
		return
	}

	// 🔹 Create new cart item
	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	// 🔹 Preload product info for response
	if err := config.DB.Preload("Product").First(&cartItem, cartItem.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Product added to cart",
		"cart_item": mapCartItem(cartItem),
	})
}

func GetCartItems(c *gin.Context) {
	userID := getUserID(c)

	var cartItems []models.CartItem
	config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems)

	var resp []CartItemResponse
	for _, item := range cartItems {
		resp = append(resp, mapCartItem(item))
	}

	c.JSON(200, gin.H{"cart_items": resp})
}

func UpdateCartItem(c *gin.Context) {
	userID := getUserID(c)
	cartID := c.Param("id")

	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var cartItem models.CartItem
	if err := config.DB.Preload("Product").Where("id = ? AND user_id = ?", cartID, userID).First(&cartItem).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cart item not found"})
		return
	}

	if input.Quantity > cartItem.Product.StockQuantity {
		c.JSON(400, gin.H{"error": "Not enough stock available"})
		return
	}

	cartItem.Quantity = input.Quantity
	config.DB.Save(&cartItem)

	c.JSON(200, gin.H{"message": "Cart updated", "cart_item": mapCartItem(cartItem)})
}

func DeleteCartItem(c *gin.Context) {
	userID := getUserID(c)
	cartID := c.Param("id")

	result := config.DB.Where("id = ? AND user_id = ?", cartID, userID).Delete(&models.CartItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart item"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}

// Helper Functions
func getUserID(c *gin.Context) uint {
	userIDValue, _ := c.Get("userID")
	switch v := userIDValue.(type) {
	case int:
		return uint(v)
	case uint:
		return v
	default:
		return 0
	}
}

func mapCartItem(item models.CartItem) CartItemResponse {
	return CartItemResponse{
		ID: item.ID,
		Product: ProductSummary{
			ID:            item.Product.ID,
			Name:          item.Product.Name,
			Description:   item.Product.Description,
			Price:         item.Product.Price,
			StockQuantity: item.Product.StockQuantity,
			ImageURL:      item.Product.ImageURL,
		},
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
