package middlewares

import (
	"IAM/initializers"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
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
		// use custom claims with user_id
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return initializers.JwtSecretKey, nil
		})
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// check token valid
		if !token.Valid {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID := claims.UserID
		c.Set("userID", userID)
		c.Next()
	}
}
