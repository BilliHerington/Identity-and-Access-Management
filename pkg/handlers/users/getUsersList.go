package users

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsersList() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get all Users from user-list in redis
		users, err := UserManageRepo.GetUsersListFromDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			if err.Error() == "redisUsers not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"users id`s": users})
	}
}
