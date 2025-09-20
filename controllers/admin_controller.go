package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
)

// UpdateUserHandler allows admin to update user info or role
func UpdateUserHandler(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		FullName  string `json:"full_name"`
		Role      string `json:"role"`
		Address   string `json:"address"`
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
	if input.Role != "" {
		updates["role"] = input.Role
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.AvatarURL != "" {
		updates["avatar_url"] = input.AvatarURL
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	updates["updated_at"] = time.Now()

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
	})
}

// BlockUserHandler allows admin to block a user
func BlockUserHandler(c *gin.Context) {
	userID := c.Param("id")

	if err := config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_blocked", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User blocked successfully",
	})
}

// UnblockUserHandler allows admin to unblock a user
func UnblockUserHandler(c *gin.Context) {
	userID := c.Param("id")

	if err := config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_blocked", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User unblocked successfully",
	})
}
