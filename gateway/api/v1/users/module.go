package users

import (
	"gateway/configs"
	"gateway/models"
	"gateway/utils"

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

	// Seed sysadmin nếu chưa có
	var count int64
	db.Model(&models.User{}).Where("role = ?", "sysadmin").Count(&count)
	if count == 0 {
		hashed, _ := utils.HashPassword("SysAdmin@123")
		db.Create(&models.User{
			Name:     "System Admin",
			Email:    "sysadmin@system.local",
			Password: hashed,
			Role:     "sysadmin",
		})
	}

	UserRepo = NewUserRepository(db)
	UserSvc = NewUserService(UserRepo)
	controller := NewUserController(UserSvc)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
