package routes

import (
	v1 "gateway/api/v1"
	"gateway/api/v1/users"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/api/v1/health", v1.HealthCheck)
	users.InitModule(r)
	return r
}
