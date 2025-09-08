package auth

import (
	"gateway/api/v1/users"
	"net/http"
	"time"

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
		authRoutes.POST("/login", c.Login)
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

func (c *AuthController) Login(ctx *gin.Context) {
	var dto LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	user, err := c.userService.ValidateUser(dto.Email, dto.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Tạo JWT token
	token, err := c.service.GenerateJWT(JWTPayload{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.SetCookie("jwt", token, int((24 * time.Hour).Seconds()), "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
