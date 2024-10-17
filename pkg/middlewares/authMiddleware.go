package middlewares

import (
	"IAM/pkg/jwt/middlewareJWT"
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from header
		tokenString, err := middlewareJWT.ExtractHeaderToken(c)
		if err != nil {
			logs.Error.Printf("No token found :%s", tokenString)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		// check token valid
		tokenValid, userData, err := middlewareJWT.CheckTokenValid(tokenString)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			return
		} else if !tokenValid {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		email, userID, jwtString := userData[0], userData[1], userData[2]
		c.Set("jwt", jwtString)
		c.Set("email", email)
		c.Set("userID", userID)
		c.Next()
	}
}
