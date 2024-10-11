package auxiliary

import (
	"IAM/pkg/logs"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestLimitRepository interface {
	GetRequestLimit(redisKey string, limit int, window int64) (bool, error)
}

func RequestLimiter(repo RequestLimitRepository, limit int, window int64) gin.HandlerFunc {
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

		// check requests exceed
		exceeded, err := repo.GetRequestLimit(redisKey, limit, window)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(500, gin.H{"error": "please try again later"})
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
