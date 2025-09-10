package users

import (
	"gateway/configs"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

// Export instance
var (
	UserRepo *UserRepository
	UserSvc  *UserService
)

func InitModule(r *gin.Engine) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.User{})
	}

	UserRepo = NewUserRepository(db)
	UserSvc = NewUserService(UserRepo)
	controller := NewUserController(UserSvc)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
