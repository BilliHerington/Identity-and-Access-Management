package roles

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type AssignRoleRepository interface {
	AssignRoleToUser(userID, role string) error
}

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

		// get userID
		userID, err := auxiliary.UserIDByEmail(&redisAuxiliaryHandlers.RedisUserIDByEmailRepo{RDB: rdb}, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// assign role
		err = AssignRoleRepository.AssignRoleToUser(&redisRolesHandlers.RedisAssignRoleRepo{RDB: rdb}, userID, input.Role)
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
