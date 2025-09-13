package models

import "time"

type OTP struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	OTPCode   string    `gorm:"type:varchar(10);not null" json:"otp_code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Purpose   string    `gorm:"type:varchar(50);not null" json:"purpose"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
