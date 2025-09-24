package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func CartRoutes(r *gin.Engine) {
	cart := r.Group("/cart") // protected routes (user only)
	cart.Use(middlewares.AuthMiddleware())
	{
		cart.POST("/", controllers.AddToCart)
		// cart.GET("/", controllers.GetCart)
		// cart.PUT("/:id", controllers.UpdateCart)
		// cart.DELETE("/:id", controllers.DeleteCart)
	}
}
