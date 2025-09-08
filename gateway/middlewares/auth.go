package middlewares

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(configs.Env.JwtSecretAccess), nil
		})
		if err != nil || !token.Valid {
			c.Error(err)
			c.Abort()
			return
		}

		// Lấy claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("name", claims["name"])
		}

		c.Next()
	}
}
