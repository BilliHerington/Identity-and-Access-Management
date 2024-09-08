package middlewares

import (
	"IAM/initializers"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
			c.Abort() // Прерываем обработку запроса
			return
		}
		// Убираем префикс "Bearer " перед токеном
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Парсим и проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем тип алгоритма шифрования
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return initializers.JwtSecretKey, nil // Возвращаем секретный ключ для проверки подписи токена
		})
		// Шаг 4: Если токен недействителен или произошла ошибка, возвращаем ошибку
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // Прерываем выполнение запроса
			return
		}
		// Шаг 5: Если токен валиден, продолжаем выполнение запроса
		c.Next()
	}
}
