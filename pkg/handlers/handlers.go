package handlers

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Define the secret key for signing JWT tokens
var jwtSecretKey = []byte("your-very-secret-key")

func generateUniqueID() string {
	return uuid.New().String()
}
func Registration(c *gin.Context) {
	var input models.RegisterData

	// 1. Получение данных от клиента и связывание с моделью User
	if err := c.ShouldBind(&input); err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. Проверка существования пользователя по email
	ctx := context.Background()
	_, err := initializers.Rdb.Get(ctx, input.Email).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logs.Error.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Пользователь уже существует
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Password = string(hashedPassword)

	// 4. Генерация уникального ID пользователя
	input.ID = generateUniqueID()

	// 5. Сериализация данных пользователя в JSON
	userData, err := json.Marshal(&input)
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 6. Сохранение пользователя в Redis
	err = initializers.Rdb.Set(ctx, input.Email, userData, 0).Err()
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 7. Отправка ответа о успешной регистрации
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
func Authenticate(c *gin.Context) {
	var input models.AuthData

	if err := c.ShouldBind(&input); err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	// Get the hashed password from Redis
	_, err := initializers.Rdb.Get(ctx, input.Email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email does not exist"})
			return
		}
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(input.Password), []byte(input.Password))
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": input.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}
