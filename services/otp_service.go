package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// GenerateOTP creates a random numeric OTP, stores it in DB with expiry
func GenerateOTP(userID uint, purpose string) (string, error) {
	// 1. Generate a random 6-digit OTP
	otp, err := generateRandomOTP(6)
	if err != nil {
		return "", err
	}

	// 2. Expiration = 5 minutes from now
	expiresAt := time.Now().Add(5 * time.Minute)

	// 3. Save to DB
	otpEntry := models.OTP{
		UserID:    userID,
		OTPCode:   otp,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&otpEntry).Error; err != nil {
		return "", err
	}

	return otp, nil
}

// ValidateOTP checks if OTP exists, not expired, and matches
func ValidateOTP(userID uint, otp string, purpose string) (bool, error) {
	var otpEntry models.OTP

	// Find OTP by userID + purpose
	err := config.DB.Where("user_id = ? AND purpose = ?", userID, purpose).
		Order("created_at DESC"). // in case multiple
		First(&otpEntry).Error
	if err != nil {
		return false, err
	}

	// Check expiry
	// Check expiry
	if time.Now().After(otpEntry.ExpiresAt) {
		return false, fmt.Errorf("otp expired")
	}

	// Check code
	if otpEntry.OTPCode != otp {
		return false, fmt.Errorf("invalid otp")
	}

	return true, nil
}

// Utility: Generate random numeric OTP
func generateRandomOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp += string(digits[n.Int64()])
	}
	return otp, nil
}
