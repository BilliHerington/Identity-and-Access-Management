package users

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type GetAllUsersDataRepository interface {
	GetAllUsersDataFromDB() (map[string]string, error)
}

func GetAllUsersData(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		repo := &redisUsersHandlers.RedisGetAllUsersDataRepo{RDB: rdb}
		data, err := repo.GetAllUsersDataFromDB()
		if errors.Is(err, redis.Nil) {
			c.Status(http.StatusNoContent)
		} else if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
