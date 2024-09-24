package roles

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AssignRole Назначение роли пользователю
func AssignRole(c *gin.Context) {
	var input models.UserRoleData
	if err := c.ShouldBindJSON(&input); err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	id, err := handlers.GetUserIDByEmail(ctx, input.Email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = initializers.Rdb.HSet(ctx, "user:"+id, "role", input.Role).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	text := fmt.Sprintf("role: %s, assign successfully for user: %s", input.Role, input.Email)
	c.JSON(http.StatusOK, gin.H{"message": text})
}
