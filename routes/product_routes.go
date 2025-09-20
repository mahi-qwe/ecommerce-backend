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
		admin.PUT("/products/:id", controllers.UpdateProductHandler)                            // PUT /admin/products/:id
		admin.DELETE("/products/:id", controllers.DeleteProductHandler)                         // DELETE /admin/products/:id
		admin.POST("/products/:id/production", controllers.StartProductionHandler)              // POST /admin/products/:id/production
		admin.PUT("/products/:id/production/status", controllers.UpdateProductionStatusHandler) // PUT /admin/products/:id/production/status
		admin.GET("/products/:id/production", controllers.GetProductionDetailsHandler)
		admin.GET("/products/production", controllers.GetAllProductionsHandler)
	}

	// Public routes
	public := r.Group("/products")
	{
		public.GET("", controllers.GetProductsHandler)        // GET /products
		public.GET("/:id", controllers.GetProductByIDHandler) // GET /products/:id
	}
}
