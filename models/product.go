package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity int            `gorm:"not null;default:0" json:"stock_quantity"`
	Category      string         `gorm:"type:varchar(100)" json:"category"`
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
