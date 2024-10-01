package auxiliary

import (
	redisDB "IAM/initializers/redisSystem"
	"IAM/pkg/logs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func RequestLimiter(limit int, window int64, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var redisKey string
		// try get userID from header
		userID := c.GetString("userID")
		if userID != "" {
			redisKey = fmt.Sprintf("rate_limit_%s", userID)
		} else {
			// if user not authorized, use IP
			clientIP := c.ClientIP()
			redisKey = fmt.Sprintf("rate_limit_ip_%s", clientIP)
		}
		exceeded, err := redisDB.RateLimitExceeded(redisKey, limit, window, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		if exceeded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
