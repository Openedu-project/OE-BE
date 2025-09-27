package wishlists

import (
	"gateway/configs"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.Wishlist{})
	}

	repo := NewWishlistRepository(db)
	service := NewWishlistService(repo)
	controller := NewWishlistController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
