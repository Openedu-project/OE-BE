package enrollments

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	repo := NewEnrollRepository(db)
	service := NewEnrollService(repo)
	controller := NewEnrollmentController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
