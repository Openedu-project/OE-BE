package auth

import (
	"gateway/api/v1/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service     *AuthService
	userService *users.UserService
}

func NewAuthController(s *AuthService, u *users.UserService) *AuthController {
	return &AuthController{
		service:     s,
		userService: u,
	}
}

func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", c.Register)
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var dto users.CreateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	user, err := c.userService.CreateUser(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
