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
	scv := NewBlogService(repo)
	controller := NewBlogController(scv)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
