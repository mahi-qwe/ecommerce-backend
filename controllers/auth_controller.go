package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/config"
	"github.com/mahi-qwe/ecommerce-backend/models"
	"github.com/mahi-qwe/ecommerce-backend/services"
	"github.com/mahi-qwe/ecommerce-backend/utils"
)

// SignupHandler handles new user registration
func SignupHandler(c *gin.Context) {
	var input struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
		return
	}

	user := models.User{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
		IsVerified:   false, // Not verified yet
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Generate OTP (already sends email inside)
	if _, err := services.GenerateOTP(user.ID, user.Email, "signup"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Account created. Please check your email for the verification code.",
	})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Block if not verified
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before login"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(int(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"userId": user.ID,
		"token":  token,
	})
}

// VerifyOTPHandler verifies the OTP sent to user's email
func VerifyOTPHandler(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		OTP     string `json:"otp" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate OTP (marks OTP used & user verified if signup)
	valid, err := services.ValidateOTP(user.ID, input.OTP, input.Purpose)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User verified successfully",
	})
}

// SendOTPHandler sends an OTP to user's email for given purpose
func SendOTPHandler(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"` // e.g., "signup", "reset_password"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate OTP (already sends email inside)
	if _, err := services.GenerateOTP(user.ID, user.Email, input.Purpose); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate/send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP sent successfully",
	})
}
