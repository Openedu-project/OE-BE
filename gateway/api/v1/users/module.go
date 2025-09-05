package users

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&User{})
	}

	repo := NewUserRepository(db)
	service := NewUserService(repo)
	controller := NewUserController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
