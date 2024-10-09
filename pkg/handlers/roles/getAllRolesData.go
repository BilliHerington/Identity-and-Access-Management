package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type GetAllRolesDataRepository interface {
	GetAllRolesDataFromDB() ([]map[string]string, error)
}

func GetAllRolesData(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := &redisRolesHandlers.RedisGetAllRolesDataRepo{RDB: rdb}
		roleData, err := repo.GetAllRolesDataFromDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			if err.Error() == "roles not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"roles": roleData})
	}
}
