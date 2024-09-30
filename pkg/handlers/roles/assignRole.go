package roles

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

// AssignRole set role for user
func AssignRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get data from client and binding with JSON
		var input models.UserRoleData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx := context.Background()
		id, err := auxiliary.GetUserIDByEmail(ctx, input.Email, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = rdb.HSet(ctx, "user:"+id, "role", input.Role).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		text := fmt.Sprintf("role: %s, assign successfully for user: %s", input.Role, input.Email)
		logs.AuditLogger.Printf(text)
		c.JSON(http.StatusOK, gin.H{"message": text})
	}
}
