package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/demkowo/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(resp.Error(http.StatusUnauthorized, "authorization header missing", nil).JSON())
			return
		}
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.AbortWithStatusJSON(resp.Error(http.StatusUnauthorized, "authorization header must start with 'Bearer '", nil).JSON())
			return
		}
		tokenStr := strings.TrimSpace(authHeader[len("bearer "):])

		jwtToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !jwtToken.Valid {
			c.AbortWithStatusJSON(resp.Error(http.StatusUnauthorized, "invalid or expired token", []interface{}{err.Error()}).JSON())
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(resp.Error(http.StatusUnauthorized, "invalid token claims", nil).JSON())
			return
		}

		idStr, ok := claims["id"].(string)
		if !ok {
			c.AbortWithStatusJSON(resp.Error(http.StatusUnauthorized, "token missing user id", nil).JSON())
			return
		}

		c.Set("account_id", idStr)

		if roles, ok := claims["roles"]; ok {
			c.Set("roles", roles)
		}

		c.Next()
	}
}
