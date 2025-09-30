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
	}

	adminPayments := r.Group("/admin/payments")
	adminPayments.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		payments.PUT("/:payment_id/update", controllers.UpdatePaymentStatus)
	}
}
