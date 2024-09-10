package access

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllUsersData(c *gin.Context) {
	ctx := context.Background()

	// Получение всех email-ов пользователей
	emails, err := initializers.Rdb.SMembers(ctx, "users").Result()
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. Для каждого email получаем данные пользователя
	var users []map[string]string
	for _, email := range emails {
		userData, err := initializers.Rdb.HGetAll(ctx, "user:"+email).Result()
		if err != nil {
			logs.Error.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(userData) > 0 {
			users = append(users, userData)
		}
	}
	// Возвращаем список пользователей
	c.JSON(http.StatusOK, users)
}
