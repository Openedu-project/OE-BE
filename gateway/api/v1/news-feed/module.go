package newsfeed

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
	}

	repo := NewBlogRepository(db)
	catRepo := NewCategoryRepository(db)
	scv := NewBlogService(repo, catRepo)
	controller := NewBlogController(scv)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
