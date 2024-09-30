package users

import (
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func GetAllUsersData(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// get all ID from user-list in redis
		IDs, err := rdb.SMembers(ctx, "users").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// get all data for every ID
		var users []map[string]string
		for _, id := range IDs {
			userData, err := rdb.HGetAll(ctx, "user:"+id).Result()
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(userData) > 0 {
				users = append(users, userData)
			}
		}
		c.JSON(http.StatusOK, users)
	}
}
