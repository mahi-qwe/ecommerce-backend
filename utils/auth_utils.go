package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ------------------ Password Hashing ------------------
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ------------------ JWT Functions ------------------
func GenerateJWT(userID int, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["userId"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid token claims")
		}
		return int(userIDFloat), nil
	}

	return 0, fmt.Errorf("invalid token")
}

// ------------------ Refresh Token Functions ------------------

// Generate a random refresh token (plain + hashed)
func GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(hash[:]), nil
}

// Save refresh token to DB
func SaveRefreshToken(db *gorm.DB, userID uint, hashedToken string, expiresAt time.Time) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}
	return db.Create(&rt).Error
}

// Validate refresh token from DB
func ValidateRefreshToken(db *gorm.DB, token string) (*models.RefreshToken, error) {
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	var rt models.RefreshToken
	err := db.Where("token = ? AND expires_at > ?", hashedToken, time.Now()).First(&rt).Error
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	return &rt, nil
}

// Delete refresh token from DB (logout)
func DeleteRefreshToken(db *gorm.DB, token string) error {
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	return db.Where("token = ?", hashedToken).Delete(&models.RefreshToken{}).Error
}
