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
		// check token valid
		tokenValid, userID, err := jwtHandlers.IsTokenValid(tokenString, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			if err.Error() == "userVersion not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
			return
		} else if !tokenValid {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
