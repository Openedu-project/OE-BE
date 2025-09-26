package enrollments

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewEnrollmentController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
