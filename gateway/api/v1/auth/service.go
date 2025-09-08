package auth

import (
	"fmt"
	"gateway/configs"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(repo *AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GenerateJWT(payload JWTPayload) (string, error) {
	secret := configs.Env.JwtSecretAccess
	expiredHours, err := strconv.Atoi(configs.Env.JwtExpiredTime)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRED_TIME: %v", err)
	}

	claims := jwt.MapClaims{
		"user_id": payload.UserID,
		"name":    payload.Name,
		"email":   payload.Email,
		"exp":     time.Now().Add(time.Duration(expiredHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
