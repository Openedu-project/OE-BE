package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})

			c.Abort()
			return
		}
	}
}
