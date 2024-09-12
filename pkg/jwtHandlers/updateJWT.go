package jwtHandlers

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJWT(c *gin.Context, id string, email string) {
	signedToken, err := CreateJWT(c, email)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	err = initializers.Rdb.HSet(ctx, "user:"+id, "jwt", signedToken).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "jwt updated successfully"})
}
