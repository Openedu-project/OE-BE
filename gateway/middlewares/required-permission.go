package middlewares

import (
	"gateway/guards"

	"github.com/gin-gonic/gin"
)

func RequirePermission(perm guards.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// role := c.GetString("role")

		// if !guards.HasPermission(guards.Role(role), perm) {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 		"error": "forbidden",
		// 	})
		// 	return
		// }

		c.Next()
	}
}
