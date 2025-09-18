package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/profile", controllers.GetProfileHandler)
		user.PUT("/profile", controllers.UpdateProfileHandler)
	}
}
