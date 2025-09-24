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
		admin.PUT("/products/:id", controllers.UpdateProductHandler)
		admin.DELETE("/products/:id", controllers.DeleteProductHandler)
		admin.POST("/products/:id/production", controllers.StartProductionHandler)              // start production route
		admin.PUT("/products/:id/production/status", controllers.UpdateProductionStatusHandler) // update production status route
		admin.GET("/products/:id/production", controllers.GetProductionDetailsHandler)          // get production details route
		admin.GET("/products/production", controllers.GetAllProductionsHandler)                 // get all productions route
	}

	// Public routes
	public := r.Group("/products")
	{
		public.GET("", controllers.GetProductsHandler)
		public.GET("/:id", controllers.GetProductByIDHandler)
	}
}
