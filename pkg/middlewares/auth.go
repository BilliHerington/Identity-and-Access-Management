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

		tokenString, err := jwtHandlers.ExtractHeaderToken(c)
		if err != nil {
			logs.Error.Printf("No token found :%s", tokenString)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			//logs.Info.Print(claims)
			return initializers.JwtSecretKey, nil
		})
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Шаг 4: Если токен недействителен или произошла ошибка, возвращаем ошибку
		if !token.Valid {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // Прерываем выполнение запроса
			return
		}
		userID := claims.UserID
		//logs.Info.Printf("userID:%s", userID)
		c.Set("userID", userID)
		// Шаг 5: Если токен валиден, продолжаем выполнение запроса
		c.Next()
	}
}
