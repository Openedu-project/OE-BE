package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *AuthService
}

func NewAuthController(s *AuthService) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("login", c.Login)
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	return
}
