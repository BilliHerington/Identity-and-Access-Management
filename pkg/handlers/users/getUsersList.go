package users

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type GetUsersListRepository interface {
	GetUsersListFromDB() ([]string, error)
}

func GetUsersList(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get all Users from user-list in redis
		repo := redisUsersHandlers.RedisGetUsersListRepo{RDB: rdb}
		users, err := repo.GetUsersListFomDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			if err.Error() == "users not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}
