package routes

import (
	v1 "gateway/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.GET("/api/v1/health", v1.HealthCheck)
	return r
}
