package services

import (
	"errors"
	"time"

	"github.com/mahi-qwe/ecommerce-backend/models"
	"gorm.io/gorm"
)

type OrderItemResponse struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	TotalAmount float64             `json:"total_amount"`
	Address     string              `json:"address"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UserName    string              `json:"user_name"`
	Items       []OrderItemResponse `json:"items"`
}

func CreateOrder(db *gorm.DB, userID uint, address string) (*OrderResponse, error) {
	var cartItems []models.CartItem
	if err := db.Where("user_id = ?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	tx := db.Begin()

	order := models.Order{
		UserID:      userID,
		Address:     address,
		Status:      "pending",
		TotalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	totalAmount := 0.0

	for _, item := range cartItems {
		if item.Product.StockQuantity < item.Quantity {
			tx.Rollback()
			return nil, errors.New("insufficient stock for product: " + item.Product.Name)
		}

		item.Product.StockQuantity -= item.Quantity
		if err := tx.Save(&item.Product).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		itemTotal := float64(item.Quantity) * item.Product.Price
		totalAmount += itemTotal

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Fetch user and order items for response
	var fullOrder models.Order
	if err := db.Preload("User").Preload("OrderItems.Product").First(&fullOrder, order.ID).Error; err != nil {
		return nil, err
	}

	// Build response DTO
	items := make([]OrderItemResponse, 0)
	for _, oi := range fullOrder.OrderItems {
		items = append(items, OrderItemResponse{
			ProductID: oi.ProductID,
			Name:      oi.Product.Name,
			Quantity:  oi.Quantity,
			Price:     oi.Price,
		})
	}

	resp := &OrderResponse{
		ID:          fullOrder.ID,
		TotalAmount: fullOrder.TotalAmount,
		Address:     fullOrder.Address,
		Status:      fullOrder.Status,
		CreatedAt:   fullOrder.CreatedAt,
		UserName:    fullOrder.User.FullName,
		Items:       items,
	}

	return resp, nil
}

// New function to get all orders for a user
func GetUserOrders(db *gorm.DB, userID uint) ([]OrderResponse, error) {
	var orders []models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	var resp []OrderResponse
	for _, order := range orders {
		items := make([]OrderItemResponse, 0)
		for _, oi := range order.OrderItems {
			items = append(items, OrderItemResponse{
				ProductID: oi.ProductID,
				Name:      oi.Product.Name,
				Quantity:  oi.Quantity,
				Price:     oi.Price,
			})
		}

		resp = append(resp, OrderResponse{
			ID:          order.ID,
			TotalAmount: order.TotalAmount,
			Address:     order.Address,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UserName:    order.User.FullName,
			Items:       items,
		})
	}

	return resp, nil
}

// Returns all orders for admin
func GetAllOrders(db *gorm.DB) ([]OrderResponse, error) {
	var orders []models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	var resp []OrderResponse
	for _, order := range orders {
		items := make([]OrderItemResponse, 0)
		for _, oi := range order.OrderItems {
			items = append(items, OrderItemResponse{
				ProductID: oi.ProductID,
				Name:      oi.Product.Name,
				Quantity:  oi.Quantity,
				Price:     oi.Price,
			})
		}

		resp = append(resp, OrderResponse{
			ID:          order.ID,
			TotalAmount: order.TotalAmount,
			Address:     order.Address,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UserName:    order.User.FullName,
			Items:       items,
		})
	}

	return resp, nil
}

// Update order status (Admin)
func UpdateOrderStatusAdmin(db *gorm.DB, orderID uint, newStatus string) (*OrderResponse, error) {
	var order models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Optional: validate status transitions
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}

	if !validStatuses[newStatus] {
		return nil, errors.New("invalid status")
	}

	order.Status = newStatus
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		return nil, err
	}

	// Build DTO response
	items := make([]OrderItemResponse, 0)
	for _, oi := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ProductID: oi.ProductID,
			Name:      oi.Product.Name,
			Quantity:  oi.Quantity,
			Price:     oi.Price,
		})
	}

	resp := &OrderResponse{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		Address:     order.Address,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UserName:    order.User.FullName,
		Items:       items,
	}

	return resp, nil
}

// Fetch single order for a specific user
func GetOrderByID(db *gorm.DB, orderID uint, userID uint) (*OrderResponse, error) {
	var order models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Security: make sure user can only access their own orders
	if order.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// Build DTO response
	items := make([]OrderItemResponse, 0)
	for _, oi := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ProductID: oi.ProductID,
			Name:      oi.Product.Name,
			Quantity:  oi.Quantity,
			Price:     oi.Price,
		})
	}

	resp := &OrderResponse{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		Address:     order.Address,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UserName:    order.User.FullName,
		Items:       items,
	}

	return resp, nil
}
