package auth

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	repo := NewAuthRepository(db)
	service := NewAuthService(repo)
	controller := NewAuthController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
