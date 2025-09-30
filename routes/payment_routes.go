package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func PaymentRoutes(r *gin.Engine) {
	payments := r.Group("/payments")
	payments.Use(middlewares.AuthMiddleware())
	{
		payments.POST("/create", controllers.CreatePaymentIntent)
		payments.PUT("/:payment_id/update", controllers.UpdatePaymentStatus)
	}
}
