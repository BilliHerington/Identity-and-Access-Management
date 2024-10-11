package jwtHandlers

import (
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func UpdateJWT(c *gin.Context, userID, userVersion, email string, rdb *redis.Client) {
	// sign token
	signedToken, err := CreateJWT(email, userVersion, rdb)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		if err.Error() == "email not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		} else {
			c.JSON(500, gin.H{"error": "please try again later"})
		}
		return
	}
	ctx := context.Background()
	// save new JWT in redis
	err = rdb.HSet(ctx, "user:"+userID, "jwt", signedToken).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(500, gin.H{"error": "please try again later"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jwt": signedToken})
}
