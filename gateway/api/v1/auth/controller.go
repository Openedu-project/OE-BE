package auth

import (
	"net/http"
	"time"

	"gateway/api/v1/users"
	wallets "gateway/api/v1/wallets"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service       *AuthService
	userService   *users.UserService
	walletService *wallets.WalletService
}

func NewAuthController(s *AuthService, u *users.UserService, w *wallets.WalletService) *AuthController {
	return &AuthController{
		service:       s,
		userService:   u,
		walletService: w,
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

	wallet, err := c.walletService.CreateNearWallet(user.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"wallet":  wallet,
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
		Role:   user.Role,
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
