package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// GenerateOTP creates, stores, and emails a 6-digit OTP
func GenerateOTP(userID uint, email, purpose string) (string, error) {
	otp, err := generateRandomOTP(6)
	if err != nil {
		return "", err
	}

	otpEntry := models.OTP{
		UserID:    userID,
		OTPCode:   otp,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
		IsUsed:    false,
	}

	if err := config.DB.Create(&otpEntry).Error; err != nil {
		return "", err
	}

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP for %s is: %s. It expires in 5 minutes.", purpose, otp)

	if err := SendEmail(email, subject, body); err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	return otp, nil
}

// ValidateOTP checks OTP validity and marks it used
func ValidateOTP(userID uint, otp, purpose string) (bool, error) {
	var entry models.OTP
	err := config.DB.Where("user_id = ? AND purpose = ? AND is_used = ?", userID, purpose, false).
		Order("created_at DESC").
		First(&entry).Error
	if err != nil {
		return false, fmt.Errorf("otp not found or already used")
	}

	if time.Now().After(entry.ExpiresAt) {
		return false, fmt.Errorf("otp expired")
	}

	if entry.OTPCode != otp {
		return false, fmt.Errorf("invalid otp")
	}

	entry.IsUsed = true
	if err := config.DB.Save(&entry).Error; err != nil {
		return false, fmt.Errorf("failed to update otp status: %w", err)
	}

	// For signup OTPs, mark user as verified
	if purpose == "signup" {
		if err := config.DB.Model(&models.User{}).
			Where("id = ?", userID).
			Update("is_verified", true).Error; err != nil {
			return false, fmt.Errorf("failed to verify user: %w", err)
		}
	}

	return true, nil
}

// generateRandomOTP returns a numeric OTP of given length
func generateRandomOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}

	return string(otp), nil
}
