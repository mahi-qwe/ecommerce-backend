package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint           `gorm:"not null" json:"order_id"`
	Gateway   string         `gorm:"type:varchar(50);not null" json:"gateway"`     // e.g., "Razorpay"
	PaymentID string         `gorm:"type:varchar(100);not null" json:"payment_id"` // ID returned by gateway
	Amount    float64        `gorm:"not null" json:"amount"`
	Status    string         `gorm:"type:varchar(50);not null" json:"status"` // success, failed, pending
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relation
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order"`
}
