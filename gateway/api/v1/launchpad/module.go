package launchpad

import (
	"gateway/configs"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

// Export instance
// var (
// 	CourseRepo *CourseRepository
// 	CourseSvc  *CourseService
// )

func InitModule(r *gin.Engine) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.Launchpad{}, &models.VotingPlan{})
	}

	repo := NewLaunchpadRepository(db)
	scv := NewLaunchpadService(repo)
	controller := NewLaunchpadController(scv)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
