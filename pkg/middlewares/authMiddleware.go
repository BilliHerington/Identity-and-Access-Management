package middlewares

import (
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from header
		tokenString, err := jwtHandlers.ExtractHeaderToken(c)
		if err != nil {
			logs.Error.Printf("No token found :%s", tokenString)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		userID := "0e631cfe"
		// check token valid
		tokenValid, err := jwtHandlers.IsTokenValid(tokenString, userID, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		} else if !tokenValid {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
