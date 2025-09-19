package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`        // Foreign key to users
	Token     string    `gorm:"not null;unique"` // Hashed token
	ExpiresAt time.Time `gorm:"not null"`        // Expiration
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
}
