package users

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
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
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if len(data) == 0 {
			c.JSON(http.StatusNoContent, gin.H{})
		}
		c.JSON(http.StatusOK, data)
	}
}
