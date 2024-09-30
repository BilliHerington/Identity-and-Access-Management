package users

import (
	"IAM/pkg/gmail"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func StartResetPassword(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// getting data from client and binding
		var input struct {
			Email string `json:"email"`
		}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check email exist
		emailMatch, err := auxiliary.EmailMatch(input.Email, rdb)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if !emailMatch {
			c.JSON(http.StatusConflict, gin.H{"error": "Email not found"})
			return
		}

		//generate code
		resetPassCode := GenerateVerificationCode()

		// email compose
		subject := "Resetting password"
		body := fmt.Sprintf("Code for resetting:%s", resetPassCode)

		// sending
		err = gmail.SendEmail(subject, body, input.Email)
		if err != nil {
			logs.ErrorLogger.Errorln(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// get userID from redis
		userID, err := auxiliary.GetUserIDByEmail(ctx, input.Email, rdb)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// add resetPassCode field to User in redis
		err = rdb.HSet(ctx, "user:"+userID, map[string]interface{}{
			"resetPassCode": resetPassCode,
		}).Err()
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		logs.AuditLogger.Printf("reset pass code sended to user: %s: %s", userID, input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Code sent"})
	}
}
func ApproveResetPassword(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// getting data from client and binding
		var input models.ResetPass
		err := c.ShouldBindJSON(&input)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check email exist in redis
		emailMatch, err := auxiliary.EmailMatch(input.Email, rdb)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !emailMatch {
			c.JSON(http.StatusConflict, gin.H{"error": "Email not found"})
			return
		}

		// get userID from redis
		userID, err := auxiliary.GetUserIDByEmail(ctx, input.Email, rdb)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// get reset pass code from redis
		code, err := rdb.HGet(ctx, "user:"+userID, "resetPassCode").Result()
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// compare codes
		if code != input.ResetPassCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Code Not Match"})
			return
		} else {
			// hashing pass
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			input.NewPassword = string(hashedPassword)

			// save new pass in redis
			err = rdb.HSet(ctx, "user:"+userID, map[string]interface{}{
				"password": input.NewPassword,
			}).Err()
			if err != nil {
				logs.ErrorLogger.Error(err)
				logs.Error.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			// delete resetPassCode field from redis
			err = rdb.HDel(ctx, "user:"+userID, "resetPassCode").Err()
			if err != nil {
				logs.ErrorLogger.Error(err)
				logs.Error.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			logs.AuditLogger.Printf("user: %s: %s reset password", userID, input.Email)
			c.JSON(http.StatusOK, gin.H{"msg": "Password updated successfully"})
		}
	}
}
