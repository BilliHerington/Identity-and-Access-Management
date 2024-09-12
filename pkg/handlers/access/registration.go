package access

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Registration(c *gin.Context) {
	var input models.RegisterData

	// 1. Получение данных от клиента и связывание с моделью User
	if err := c.ShouldBind(&input); err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Проверка существования пользователя по email
	emailMatch, err := handlers.EmailMatch(input.Email)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if emailMatch {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// 3. Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.Password = string(hashedPassword)

	// 4. Генерация уникального ID пользователя
	userID := uuid.New().String()[:8]

	// 5. Сохранение пользователя в Redis с помощью HSet
	ctx := context.Background()
	err = initializers.Rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
				"id":       userID,
				"email":    input.Email,
				"name":     input.Name,
				"password": input.Password,
				"role":     input.Role,
				"jwt":      input.JWT,
			})
			pipe.SAdd(ctx, "users", userID)
			return nil
		})
		return err
	}, "user:"+userID)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//6. Добавление email в отдельный ключ для связи email -> ID пользователя
	err = initializers.Rdb.Set(ctx, "email:"+input.Email, userID, 0).Err()
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 7. Отправка успешного ответа
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
