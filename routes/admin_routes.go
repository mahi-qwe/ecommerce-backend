package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware()) // protect with JWT + admin role
	{
		admin.PUT("/users/:id", controllers.UpdateUserHandler)
		admin.POST("/users/:id/block", controllers.BlockUserHandler)
		admin.POST("/users/:id/unblock", controllers.UnblockUserHandler)
	}
}
