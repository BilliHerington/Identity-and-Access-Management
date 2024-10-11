package middlewares

import (
	"IAM/pkg/logs"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func CheckPrivileges(requiredPrivilege string, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get userID from header
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		ctx := context.Background()
		// get user role from redis
		role, err := rdb.HGet(ctx, "user:"+userID, "role").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			c.Abort()
			return
		}
		// get privileges by role
		data, err := rdb.HGet(ctx, "role:"+role, "privileges").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			c.Abort()
			return
		}
		// unmarshal privileges
		var privileges []string
		err = json.Unmarshal([]byte(data), &privileges)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			c.Abort()
			return
		}
		// compare user privileges with required privileges
		hasPrivileges := false
		for _, privilege := range privileges {
			if privilege == requiredPrivilege {
				hasPrivileges = true
				break
			}
		}
		// if not match
		if !hasPrivileges {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required privileges"})
			c.Abort()
			return
		}
		c.Next()
	}
}
