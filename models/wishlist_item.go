package models

import "time"

type WishlistItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
