package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type GetRolesListRepository interface {
	GetRolesListFromDB() ([]string, error)
}

func GetRolesList(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := &redisRolesHandlers.GetRolesListRepo{RDB: rdb}
		roles, err := repo.GetRolesListFromDB()
		if errors.Is(err, redis.Nil) {
			c.Status(http.StatusNoContent)
		} else if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"roles list": roles})
	}
}
