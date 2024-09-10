package roles

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllRolesData(c *gin.Context) {
	ctx := context.Background()
	roleNames, err := initializers.Rdb.SMembers(ctx, "roles").Result()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var roles []map[string]string
	for _, roleName := range roleNames {
		roleData, err := initializers.Rdb.HGetAll(ctx, "role:"+roleName).Result()
		if err != nil {
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(roleData) > 0 {
			roles = append(roles, roleData)
		}
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
