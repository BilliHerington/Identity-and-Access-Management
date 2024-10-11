package users

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func GetAllUsersData(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := UserManageRepo.GetAllUsersDataFromDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		} else if len(data) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"error": "no data found"})
		}
		c.JSON(http.StatusOK, data)
	}
}
