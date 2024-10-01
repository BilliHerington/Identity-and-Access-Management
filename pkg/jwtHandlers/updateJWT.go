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
	signedToken, err := CreateJWT(c, email, userVersion, rdb)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	// save new JWT in redis
	err = rdb.HSet(ctx, "user:"+userID, "jwt", signedToken).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jwt": signedToken})
}
