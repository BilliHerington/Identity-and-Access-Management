package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type DeleteRoleRepository interface {
	DeleteRole(roleName string) error
}

func DeleteRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input struct {
			RoleName string `json:"role_name"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Print(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// delete role
		repo := &redisRolesHandlers.RedisDeleteRoleRepo{RDB: rdb}
		err := repo.DeleteRole(input.RoleName)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logs.AuditLogger.Printf("role deleted successfully: %s", input.RoleName)
		c.JSON(http.StatusOK, gin.H{"role deleted successfully: ": input.RoleName})
	}
}
