package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type RedactRoleRepository interface {
	RedactRoleDB(role string, privileges []string) error
}

func RedactRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.RolesData
		// get data from client and binding with JSON
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// redact role
		err := redisRolesHandlers.RedisRedactRoleRepo{RDB: rdb}.RedactRoleDB(input.RoleName, input.Privileges)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs.AuditLogger.Printf("%s updated successfully. New privileges: %s", input.RoleName, input.Privileges)
		c.JSON(http.StatusOK, gin.H{"role:" + input.RoleName + " updated successfully. New privileges": input.Privileges})
	}
}
