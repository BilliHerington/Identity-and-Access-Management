package access

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteUser(c *gin.Context) {
	//userID := c.Param("id")
	var input models.DeleteUserData
	rdb := initializers.Rdb
	ctx := context.Background()
	err := c.ShouldBindJSON(&input)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := input.UserID

	email, err := rdb.HGet(ctx, "user:"+userID, "email").Result()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = rdb.Del(ctx, "user:"+userID).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = rdb.Del(ctx, "email:"+email).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = rdb.SRem(ctx, "users", userID).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
