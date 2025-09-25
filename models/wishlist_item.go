package models

import "time"

type WishlistItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"-"`
	ProductID uint      `gorm:"not null" json:"-"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
