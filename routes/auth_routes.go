package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.SignupHandler)
		// next: login, send-otp, verify-otp will be here
		auth.POST("/login", controllers.LoginHandler)
	}
}
