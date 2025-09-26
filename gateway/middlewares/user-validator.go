package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserValidatorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdValue, exists := ctx.Get("userId")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: userId not found in context",
			})
			return
		}
		if _, ok := userIdValue.(uint); !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid userId type in context",
			})
			return
		}
		ctx.Next()
	}
}
