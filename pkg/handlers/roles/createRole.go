package roles

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

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

		roleKey := "role:" + input.RoleName
		//check role exist in redis
		match, err := auxiliary.RoleMatch(roleKey, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if match {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role already exists"})
			return
		}
		// marshal Privileges for writing in redis
		privilegesJSON, err := json.Marshal(input.Privileges)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// writing in redis
		ctx := context.Background()
		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, roleKey, map[string]interface{}{
					"name":       input.RoleName,
					"privileges": privilegesJSON,
				})
				pipe.SAdd(ctx, "roles", input.RoleName)
				return nil
			})
			return err
		}, roleKey)
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
