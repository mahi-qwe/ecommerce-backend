package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware ensures only admins can access certain routes
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role") // role should be set inside JWT claims
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
