package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/services"
)

type PlaceOrderRequest struct {
	Address string `json:"address" binding:"required"`
}

func PlaceOrder(c *gin.Context) {
	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDInt, ok := uid.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	order, err := services.CreateOrder(config.DB, uint(userIDInt), req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GET /orders - user order history
func GetUserOrders(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDInt, ok := uid.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	orders, err := services.GetUserOrders(config.DB, uint(userIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GET /admin/orders - Admin view all orders
func GetAllOrders(c *gin.Context) {
	orders, err := services.GetAllOrders(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// PUT /admin/order/:id - Admin updates order status
func UpdateOrderStatusAdmin(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrder, err := services.UpdateOrderStatusAdmin(config.DB, uint(orderID), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedOrder)
}

// GET /order/:id - Get details of a specific order for logged-in user
func GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt, ok := uid.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}
	userID := uint(userIDInt)

	order, err := services.GetOrderByID(config.DB, uint(orderID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DELETE /order/:id - Soft delete a user's order
func DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt, ok := uid.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}
	userID := uint(userIDInt)

	if err := services.DeleteOrder(config.DB, uint(orderID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}
