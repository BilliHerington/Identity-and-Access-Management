package roles

import (
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func DeleteRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
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
		role := input.RoleName
		// deleting role from role-list in redis
		err := rdb.SRem(ctx, "roles", role).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// deleting role data from redis
		err = rdb.Del(ctx, "role:"+role).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs.AuditLogger.Printf("role deleted successfully: %s", role)
		c.JSON(http.StatusOK, gin.H{"role deleted successfully: ": role})
	}
}
