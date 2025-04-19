package serviceauth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServiceAuthMiddleware(redis RedisClient, bootstrapToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := c.GetHeader("X-Service-Name")
		key := c.GetHeader("X-API-Key")

		if service == "" || key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing service name or key"})
			return
		}

		expectedKey, err := redis.GetServiceKey(service)
		if err == nil && expectedKey == key {
			c.Set("caller_service", service)
			c.Next()
			return
		}

		// warunkowy bootstrap bypass tylko dla `auth`
		if service == "auth" && bootstrapToken != "" {
			bootstrap := c.GetHeader("X-Bootstrap-Token")
			hasKey, err := redis.HasServiceKey("auth")
			if err == nil && !hasKey && bootstrap == bootstrapToken {
				log.Println("Bootstrap token accepted for first-time sync from auth")
				c.Set("caller_service", service)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
	}
}
