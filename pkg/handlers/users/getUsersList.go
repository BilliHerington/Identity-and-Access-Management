package users

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func GetUsersList(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get all Users from user-list in redis
		users, err := UserManageRepo.GetUsersListFomDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			if err.Error() == "users not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}
