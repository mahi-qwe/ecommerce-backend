package seeders

import (
	"log"

	"github.com/mahi-qwe/ecommerce-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		// Original ones
		{
			FullName:     "John Doe",
			Email:        "john@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsBlocked:    false,
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=2"),
			Address:      "45 Elm Street, Springfield",
		},
		{
			FullName:     "Jane Smith",
			Email:        "jane@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsBlocked:    false,
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=3"),
			Address:      "99 Maple Avenue, Riverdale",
		},
		// Extra 7 fake customers
		{
			FullName:     "Michael Johnson",
			Email:        "michael@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=4"),
			Address:      "12 Oak Street, Gotham",
		},
		{
			FullName:     "Emily Davis",
			Email:        "emily@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=5"),
			Address:      "77 Sunset Blvd, Los Angeles",
		},
		{
			FullName:     "Chris Brown",
			Email:        "chris@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=6"),
			Address:      "21 Main Road, Star City",
		},
		{
			FullName:     "Sophia Wilson",
			Email:        "sophia@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=7"),
			Address:      "89 Lake View, Central City",
		},
		{
			FullName:     "David Miller",
			Email:        "david@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=8"),
			Address:      "56 Pine Street, Smallville",
		},
		{
			FullName:     "Olivia Taylor",
			Email:        "olivia@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=9"),
			Address:      "14 Park Avenue, Coast City",
		},
		{
			FullName:     "James Anderson",
			Email:        "james@example.com",
			PasswordHash: string(password),
			Role:         "user",
			IsVerified:   true,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?img=10"),
			Address:      "5 Broadway Street, Metropolis",
		},
	}

	for _, u := range users {
		err := db.Where(models.User{Email: u.Email}).FirstOrCreate(&u).Error
		if err != nil {
			log.Printf("❌ Could not seed user %s: %v", u.Email, err)
		}
	}

	log.Println("✅ Users seeded")
}

func strPtr(s string) *string {
	return &s
}
