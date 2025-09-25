package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahi-qwe/ecommerce-backend/controllers"
	"github.com/mahi-qwe/ecommerce-backend/middlewares"
)

func WishlistRoutes(r *gin.Engine) {
	wishlist := r.Group("/wishlist")
	wishlist.Use(middlewares.AuthMiddleware())
	{
		wishlist.POST("", controllers.AddToWishlist)
		wishlist.GET("", controllers.GetWishlist)
		wishlist.DELETE("/:product_id", controllers.RemoveFromWishlist)
	}
}
