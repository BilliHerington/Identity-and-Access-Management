package users

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

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
		// get id
		ctx := context.Background()
		userID, err := auxiliary.GetUserIDByEmail(ctx, input.Email, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// deleting userdata from redis
		err = rdb.Del(ctx, "user:"+userID).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// deleting email from redis
		err = rdb.Del(ctx, "email:"+input.Email).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// deleting userID from list in redis
		err = rdb.SRem(ctx, "users", userID).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user deleted": userID + ": " + input.Email})
	}
}
