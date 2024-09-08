package access

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Authenticate(c *gin.Context) {
	var input models.AuthData

	// 1. Получение данных от клиента и связывание с моделью User
	if err := c.ShouldBind(&input); err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emailMatch, err := handlers.EmailMatch(c, input.Email)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !emailMatch {
		c.JSON(http.StatusConflict, gin.H{"error": "Email does not match"})
		return
	}
	id, err := handlers.GetUserIDByEmail(c, input.Email)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pass, err := initializers.Rdb.HGet(c, "user:"+id, "password").Result()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(input.Password))
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": input.Email,                           // Добавляем в токен данные (в данном случае email)
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Устанавливаем время жизни токена
	})
	// Подписываем токен с использованием секретного ключа
	tokenString, err := token.SignedString(initializers.JwtSecretKey)
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	err = initializers.Rdb.HSet(ctx, "user:"+id, "jwt", tokenString).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "JWT updated successfully"})
}
