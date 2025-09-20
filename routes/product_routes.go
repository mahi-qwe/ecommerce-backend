package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func ProductRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware()) // protect admin routes
	{
		admin.POST("/products", controllers.CreateProductHandler)
	}

	// Public routes
	public := r.Group("/products")
	{
		public.GET("", controllers.GetProductsHandler)        // GET /products
		public.GET("/:id", controllers.GetProductByIDHandler) // GET /products/:id
	}
}
