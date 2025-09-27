package seeders

import (
	"log"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	log.Println("Starting database seeding...")

	SeedUsers(db)
	SeedProducts(db)
	SeedOrders(db)

	log.Println("Database seeding complete!")
}
