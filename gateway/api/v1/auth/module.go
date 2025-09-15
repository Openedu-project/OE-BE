package auth

import (
	"gateway/api/v1/users"
	wallets "gateway/api/v1/wallets"
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	repo := NewAuthRepository(db)
	service := NewAuthService(repo)
	controller := NewAuthController(service, users.UserSvc, wallets.WalletSvc)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
