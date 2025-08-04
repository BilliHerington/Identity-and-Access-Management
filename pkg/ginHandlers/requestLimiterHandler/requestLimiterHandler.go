package requestLimiterHandler

import (
	"IAM/pkg/models"
	"IAM/pkg/service/requestLimiterService"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestLimiter(limit int, window int64) gin.HandlerFunc {
	return func(c *gin.Context) {

		// try get userID from header
		userID := c.GetString("userID")
		clientIP := c.ClientIP()

		if err := requestLimiterService.RequestLimiterService(userID, clientIP, limit, window); err != nil {

			if errors.Is(err, models.ErrRequestLimitExceeded) {
				c.AbortWithStatusJSON(http.StatusRequestTimeout, models.ErrRequestLimitExceeded)
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			c.Abort()
			return
		}
		c.Next()
	}
}
