package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"github.com/mahi-qwe/ecommerce-backend/utils"
)

// GetProfileHandler returns the logged-in user's profile
func GetProfileHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(int)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Only return safe fields
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"full_name":   user.FullName,
			"email":       user.Email,
			"is_verified": user.IsVerified,
		},
	})
}

// UpdateProfileHandler updates logged-in user's profile
func UpdateProfileHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(int)

	// Only allow safe fields
	var input struct {
		FullName  string `json:"full_name"`
		Password  string `json:"password"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.FullName != "" {
		updates["full_name"] = input.FullName
	}

	if input.AvatarURL != "" {
		updates["avatar_url"] = input.AvatarURL
	}

	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password_hash"] = hashedPassword
	}

	// If nothing to update
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	// Update timestamp
	updates["updated_at"] = time.Now()

	if err := config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
	})
}
