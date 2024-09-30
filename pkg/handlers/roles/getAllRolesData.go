package roles

import (
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func GetAllRolesData(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		roleNames, err := rdb.SMembers(ctx, "roles").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var roles []map[string]string
		for _, roleName := range roleNames {
			roleData, err := rdb.HGetAll(ctx, "role:"+roleName).Result()
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(roleData) > 0 {
				roles = append(roles, roleData)
			}
		}
		c.JSON(http.StatusOK, gin.H{"roles": roles})
	}
}
