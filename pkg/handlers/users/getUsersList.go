package users

import (
	"IAM/pkg/logs"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func GetUsersList(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// get all Users from user-list in redis
		users, err := rdb.SMembers(ctx, "users").Result()
		if errors.Is(err, redis.Nil) {
			c.Status(http.StatusNoContent)
		} else if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}
