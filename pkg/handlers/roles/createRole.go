package roles

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type CreateRoleRepository interface {
	CreateRole(roleName string, privileges []string) error
}

func CreateRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.RolesData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// transfer client Redis (from main) in type RedisRoleRepository struct
		repo := &redisAuxiliaryHandlers.RedisRoleRepo{RDB: rdb} // RDB - field of RedisRoleRepository struct, rdb - redis client from main

		//check role exist
		match, err := auxiliary.RoleMatch(repo, input.RoleName)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if match {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role already exists"})
			return
		}

		// save role
		err = redisRolesHandlers.RedisCreateRoleRepo{RDB: rdb}.CreateRole(input.RoleName, input.Privileges)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs.AuditLogger.Printf("role created successfully %s", input.RoleName)
		c.JSON(http.StatusOK, gin.H{"role created successfully": input.RoleName})
	}
}
