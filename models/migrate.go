package models

import (
	"log"

	"github.com/mahi-qwe/ecommerce-backend/config"
)

func Migrate() {
	err := config.DB.AutoMigrate(
		&User{},
		&OTP{},
		&Product{},
		&ProductProduction{},
		&CartItem{},
		&WishlistItem{},
		&Order{},
		&OrderItem{},
		&RefreshToken{},
	)

	if err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}

	log.Println("✅ All tables migrated successfully")
}
