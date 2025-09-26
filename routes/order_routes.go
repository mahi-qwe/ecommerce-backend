package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func OrderRoutes(r *gin.Engine) {
	order := r.Group("/order")
	order.Use(middlewares.AuthMiddleware())
	{
		order.POST("", controllers.PlaceOrder) // place an order
		order.GET("", controllers.GetUserOrders)
	}
}
