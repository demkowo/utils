package serviceauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusHandler(redis RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, _ := redis.HasServiceKey("auth")
		status := "not received"
		if ok {
			status = "received"
		}
		c.JSON(http.StatusOK, gin.H{
			"auth_key": status,
		})
	}
}
