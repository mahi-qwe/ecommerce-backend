package seeders

import (
	"log"

	"github.com/mahi-qwe/ecommerce-backend/models"
	"gorm.io/gorm"
)

func SeedOrders(db *gorm.DB) {
	// Fetch all active users and products
	var users []models.User
	var products []models.Product

	if err := db.Where("deleted_at IS NULL").Find(&users).Error; err != nil || len(users) == 0 {
		log.Println("❌ No active users found, cannot seed orders")
		return
	}

	if err := db.Where("deleted_at IS NULL").Find(&products).Error; err != nil || len(products) == 0 {
		log.Println("❌ No active products found, cannot seed orders")
		return
	}

	// Ensure we have at least 10 users and products
	if len(users) < 10 || len(products) < 12 {
		log.Println("⚠️ Not enough users/products for hardcoded orders")
	}

	// Hardcoded order structure (user index & product index instead of IDs)
	orderTemplates := []struct {
		UserIndex int
		Address   string
		Status    string
		Items     []struct {
			ProductIndex int
			Quantity     int
		}
	}{
		{0, "Kavanur", "pending", []struct{ ProductIndex, Quantity int }{{0, 1}, {1, 1}}},
		{1, "Chennai", "shipped", []struct{ ProductIndex, Quantity int }{{2, 1}, {3, 1}}},
		{2, "Bangalore", "delivered", []struct{ ProductIndex, Quantity int }{{4, 2}}},
		{3, "Mumbai", "processing", []struct{ ProductIndex, Quantity int }{{5, 1}}},
		{4, "Delhi", "pending", []struct{ ProductIndex, Quantity int }{{6, 2}}},
		{5, "Hyderabad", "shipped", []struct{ ProductIndex, Quantity int }{{7, 3}}},
		{6, "Pune", "delivered", []struct{ ProductIndex, Quantity int }{{8, 1}}},
		{7, "Kolkata", "pending", []struct{ ProductIndex, Quantity int }{{9, 2}}},
		{8, "Chandigarh", "processing", []struct{ ProductIndex, Quantity int }{{10, 1}}},
		{9, "Jaipur", "shipped", []struct{ ProductIndex, Quantity int }{{11, 2}}},
	}

	for _, ot := range orderTemplates {
		if ot.UserIndex >= len(users) {
			log.Printf("⚠️ User index %d out of range, skipping order", ot.UserIndex)
			continue
		}
		user := users[ot.UserIndex]

		// Check if order already exists
		var existing models.Order
		if err := db.Where("user_id = ? AND address = ?", user.ID, ot.Address).First(&existing).Error; err == nil {
			continue // skip if already exists
		}

		var orderItems []models.OrderItem
		totalAmount := 0.0
		canCreate := true

		for _, item := range ot.Items {
			if item.ProductIndex >= len(products) {
				log.Printf("⚠️ Product index %d out of range, skipping order for user %d", item.ProductIndex, user.ID)
				canCreate = false
				break
			}
			product := products[item.ProductIndex]

			if product.StockQuantity < item.Quantity {
				log.Printf("⚠️ Not enough stock for product %s, skipping order for user %d", product.Name, user.ID)
				canCreate = false
				break
			}

			// Deduct stock
			product.StockQuantity -= item.Quantity
			if err := db.Save(&product).Error; err != nil {
				log.Printf("❌ Could not update stock for product %s: %v", product.Name, err)
				canCreate = false
				break
			}

			orderItems = append(orderItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				Price:     product.Price * float64(item.Quantity),
			})
			totalAmount += product.Price * float64(item.Quantity)
		}

		if !canCreate || len(orderItems) == 0 {
			continue
		}

		order := models.Order{
			UserID:      user.ID,
			TotalAmount: totalAmount,
			Address:     ot.Address,
			Status:      ot.Status,
			OrderItems:  orderItems,
		}

		if err := db.Create(&order).Error; err != nil {
			log.Printf("❌ Could not create order for user %d: %v", user.ID, err)
		}
	}

	log.Println("✅ Hardcoded orders seeded successfully with stock updated")
}
