package middlewares

import (
	"gateway/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			status := http.StatusInternalServerError

			if appErr, ok := err.(*utils.AppError); ok {
				status = appErr.Status
			}

			c.JSON(status, ErrorResponse{
				Status:    status,
				Message:   err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request.URL.Path,
			})
			c.Abort()
			return
		}
	}
}
