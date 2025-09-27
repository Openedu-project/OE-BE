package certificates

import (
	"gateway/configs"
	"gateway/models"

	"github.com/gin-gonic/gin"
)

func InitModule(r *gin.Engine) {
	db := configs.DB

	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.Certificate{})
	}

	repo := NewCertificateRepository(db)
	service := NewCertificateService(repo)
	controller := NewCertificateController(service)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}
