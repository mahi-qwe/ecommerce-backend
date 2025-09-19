package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.SignupHandler)
		auth.POST("/login", controllers.LoginHandler)
		auth.POST("/send-otp", controllers.SendOTPHandler)
		auth.POST("/verify-otp", controllers.VerifyOTPHandler)
		auth.POST("/forgot-password", controllers.ForgotPasswordHandler)
		auth.POST("/reset-password", controllers.ResetPasswordHandler)

		// New refresh token endpoints
		auth.POST("/refresh", controllers.RefreshTokenHandler)
		auth.POST("/logout", controllers.LogoutHandler)
	}
}
