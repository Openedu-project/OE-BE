package middlewares

import (
	"strconv"

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

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			var userId uint
			switch v := claims["user_id"].(type) {
			case float64:
				userId = uint(v)
			case string:
				if id, err := strconv.ParseUint(v, 10, 32); err == nil {
					userId = uint(id)
				}
			case int:
				userId = uint(v)
			}

			c.Set("userId", userId)
			c.Set("email", claims["email"])
			c.Set("name", claims["name"])
		}

		c.Next()
	}
}
