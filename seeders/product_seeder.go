package seeders

import (
	"log"

	"github.com/mahi-qwe/ecommerce-backend/models"
	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) {
	products := []models.Product{
		// Original ones
		{
			Name:          "iPhone 15 Pro",
			Description:   "Latest Apple flagship smartphone with A17 Bionic chip.",
			Price:         129999.00,
			StockQuantity: 10,
			Category:      "Electronics",
			ImageURL:      "https://example.com/images/iphone15pro.jpg",
		},
		{
			Name:          "Samsung Galaxy S23 Ultra",
			Description:   "High-end Samsung smartphone with advanced camera features.",
			Price:         119999.00,
			StockQuantity: 15,
			Category:      "Electronics",
			ImageURL:      "https://example.com/images/galaxys23ultra.jpg",
		},
		{
			Name:          "Sony WH-1000XM5",
			Description:   "Noise-cancelling wireless headphones with premium sound.",
			Price:         29999.00,
			StockQuantity: 25,
			Category:      "Audio",
			ImageURL:      "https://example.com/images/sonywh1000xm5.jpg",
		},
		{
			Name:          "MacBook Pro 16",
			Description:   "Apple M2 Max laptop with stunning Retina display.",
			Price:         249999.00,
			StockQuantity: 5,
			Category:      "Computers",
			ImageURL:      "https://example.com/images/macbookpro16.jpg",
		},
		// Extra 7 fake products
		{
			Name:          "Nike Air Max 270",
			Description:   "Stylish and comfortable running shoes.",
			Price:         12999.00,
			StockQuantity: 40,
			Category:      "Footwear",
			ImageURL:      "https://example.com/images/nikeairmax270.jpg",
		},
		{
			Name:          "Adidas Ultraboost",
			Description:   "Premium cushioned shoes for everyday comfort.",
			Price:         13999.00,
			StockQuantity: 30,
			Category:      "Footwear",
			ImageURL:      "https://example.com/images/adidasultraboost.jpg",
		},
		{
			Name:          "Apple Watch Series 9",
			Description:   "Smartwatch with fitness tracking and health monitoring.",
			Price:         45999.00,
			StockQuantity: 20,
			Category:      "Wearables",
			ImageURL:      "https://example.com/images/applewatch9.jpg",
		},
		{
			Name:          "Dell XPS 13",
			Description:   "Compact ultrabook with powerful performance.",
			Price:         149999.00,
			StockQuantity: 12,
			Category:      "Computers",
			ImageURL:      "https://example.com/images/dellxps13.jpg",
		},
		{
			Name:          "LG OLED C2 TV",
			Description:   "55-inch 4K OLED TV with Dolby Vision & Atmos.",
			Price:         99999.00,
			StockQuantity: 8,
			Category:      "Home Entertainment",
			ImageURL:      "https://example.com/images/lgc2tv.jpg",
		},
		{
			Name:          "Logitech MX Master 3",
			Description:   "Ergonomic wireless mouse with fast scrolling.",
			Price:         7999.00,
			StockQuantity: 50,
			Category:      "Accessories",
			ImageURL:      "https://example.com/images/mxmaster3.jpg",
		},
		{
			Name:          "Canon EOS R6",
			Description:   "Mirrorless camera with 20MP full-frame sensor.",
			Price:         189999.00,
			StockQuantity: 6,
			Category:      "Cameras",
			ImageURL:      "https://example.com/images/canoneosr6.jpg",
		},
	}

	for _, p := range products {
		err := db.Where(models.Product{Name: p.Name}).FirstOrCreate(&p).Error
		if err != nil {
			log.Printf("❌ Could not seed product %s: %v", p.Name, err)
		}
	}

	log.Println("✅ Products seeded")
}
