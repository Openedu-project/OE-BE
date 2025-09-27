package certificates

import (
	"net/http"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
)

type CertificateController struct {
	service *CertificateService
}

func NewCertificateController(service *CertificateService) *CertificateController {
	return &CertificateController{service: service}
}

func (c *CertificateController) RegisterRoutes(r *gin.RouterGroup) {
	certRoutes := r.Group("/my-certificates")

	certRoutes.Use(middlewares.AuthMiddleware())
	certRoutes.Use(middlewares.UserValidatorMiddleware())
	certRoutes.Use(middlewares.RequirePermission(guards.PermViewMyCertificates))
	{
		certRoutes.GET("/", c.GetMyCertificates)
	}
}

func (c *CertificateController) GetMyCertificates(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(uint)

	certificates, err := c.service.GetMyCertificates(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve certificates",
		})
		return
	}

	ctx.JSON(http.StatusOK, certificates)
}
