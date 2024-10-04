package users

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type DeleteUserRepository interface {
	DeleteUserDB(userID, email string) error
}

func DeleteUser(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input struct {
			Email string `json:"email"`
		}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// get userID
		userID, err := auxiliary.UserIDByEmail(&redisAuxiliaryHandlers.RedisUserIDByEmailRepo{RDB: rdb}, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// deleting userdata from redis
		repo := &redisUsersHandlers.RedisDeleteUserRepo{RDB: rdb}
		err = repo.DeleteUserDB(userID, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user deleted": userID + ": " + input.Email})
	}
}
