package access

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsersList(c *gin.Context) {
	ctx := context.Background()

	// Получение всех email-ов пользователей
	roles, err := initializers.Rdb.SMembers(ctx, "roles").Result()
	if err != nil {
		logs.Error.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
