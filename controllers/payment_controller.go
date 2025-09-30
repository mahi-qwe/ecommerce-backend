package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"github.com/mahi-qwe/ecommerce-backend/services"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

// DTOs
// type OrderItemResponse struct {
// 	ProductID uint    `json:"product_id"`
// 	Name      string  `json:"name"`
// 	Quantity  int     `json:"quantity"`
// 	Price     float64 `json:"price"`
// }

// type OrderResponse struct {
// 	ID          uint                `json:"id"`
// 	TotalAmount float64             `json:"total_amount"`
// 	Address     string              `json:"address"`
// 	Status      string              `json:"status"`
// 	CreatedAt   time.Time           `json:"created_at"`
// 	UserName    string              `json:"user_name"`
// 	Items       []OrderItemResponse `json:"items"`
// }

type PaymentResponse struct {
	ID        uint      `json:"id"`
	PaymentID string    `json:"payment_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Gateway   string    `json:"gateway"`
	OrderID   uint      `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// POST /payments/create
func CreatePaymentIntent(c *gin.Context) {
	type RequestBody struct {
		OrderID uint `json:"order_id"`
	}
	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Load order with items and user
	var order models.Order
	if err := config.DB.Preload("OrderItems.Product").Preload("User").First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already initiated or order not pending"})
		return
	}

	// Check for existing pending payment
	var existingPayment models.Payment
	if err := config.DB.First(&existingPayment, "order_id = ? AND status = ?", order.ID, "pending").Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pending payment already exists"})
		return
	}

	// Set Stripe API key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Create Stripe PaymentIntent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.TotalAmount * 100)),
		Currency: stripe.String(string(stripe.CurrencyINR)),
	}
	params.AddMetadata("order_id", strconv.Itoa(int(order.ID)))

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save Payment in DB
	payment := models.Payment{
		OrderID:   order.ID,
		Gateway:   "Stripe",
		PaymentID: pi.ID,
		Amount:    order.TotalAmount,
		Status:    "pending",
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	paymentResp := PaymentResponse{
		ID:        payment.ID,
		PaymentID: payment.PaymentID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		Gateway:   payment.Gateway,
		OrderID:   payment.OrderID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"payment":       paymentResp,
		"client_secret": pi.ClientSecret,
	})
}

// PUT /payments/:payment_id/update
func UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("payment_id")
	type Body struct {
		Status string `json:"status"` // "succeeded" or "failed"
	}
	var body Body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var payment models.Payment
	if err := config.DB.Preload("Order.OrderItems.Product").Preload("Order.User").First(&payment, "payment_id = ?", paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already processed"})
		return
	}

	if body.Status != "succeeded" && body.Status != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Update payment status
	payment.Status = body.Status
	if err := config.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	// Handle order & stock
	switch body.Status {
	case "succeeded":
		for _, item := range payment.Order.OrderItems {
			if item.Product.StockQuantity >= item.Quantity {
				item.Product.StockQuantity -= item.Quantity
				config.DB.Save(&item.Product)
			}
		}
		// Clear user's cart
		config.DB.Where("user_id = ?", payment.Order.UserID).Delete(&models.CartItem{})
		payment.Order.Status = "processing"
		config.DB.Save(&payment.Order)

	case "failed":
		payment.Order.Status = "failed"
		config.DB.Save(&payment.Order)
	}

	// Build response DTOs
	orderItemsResp := []services.OrderItemResponse{}
	for _, item := range payment.Order.OrderItems {
		orderItemsResp = append(orderItemsResp, services.OrderItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	orderResp := services.OrderResponse{
		ID:          payment.Order.ID,
		TotalAmount: payment.Order.TotalAmount,
		Address:     payment.Order.Address,
		Status:      payment.Order.Status,
		CreatedAt:   payment.Order.CreatedAt,
		UserName:    payment.Order.User.FullName,
		Items:       orderItemsResp,
	}

	paymentResp := PaymentResponse{
		ID:        payment.ID,
		PaymentID: payment.PaymentID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		Gateway:   payment.Gateway,
		OrderID:   payment.OrderID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentResp,
		"order":   orderResp,
	})
}
