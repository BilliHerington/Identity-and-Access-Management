package middlewares

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckPrivileges(requiredPrivilege string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		//logs.Info.Println("userID :" + userID)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		ctx := context.Background()
		role, err := initializers.Rdb.HGet(ctx, "user:"+userID, "role").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		data, err := initializers.Rdb.HGet(ctx, "role:"+role, "privileges").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		var privileges []string
		err = json.Unmarshal([]byte(data), &privileges)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		hasPrivileges := false
		for _, privilege := range privileges {
			if privilege == requiredPrivilege {
				hasPrivileges = true
				break
			}
		}
		if !hasPrivileges {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required privileges"})
			c.Abort()
			return
		}
		c.Next()
	}
}
