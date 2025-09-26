package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	TotalAmount float64        `gorm:"not null" json:"total_amount"`
	Address     string         `gorm:"type:varchar(255);not null" json:"address"`
	Status      string         `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
}
