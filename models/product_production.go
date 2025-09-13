package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductProduction struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint           `gorm:"not null" json:"product_id"`
	Status      string         `gorm:"type:varchar(50);not null" json:"status"` // pending, in_progress, completed
	StartedAt   time.Time      `gorm:"autoCreateTime" json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relation
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
