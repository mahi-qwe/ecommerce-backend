package models

import (
	"time"
)

type Order struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	TotalAmount float64    `gorm:"not null" json:"total_amount"`
	Status      string     `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"` // soft delete

	// Relation: an order belongs to a user
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"` // relation to order items
}
