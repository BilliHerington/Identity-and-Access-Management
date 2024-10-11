package users

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/handlers/gmail"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func StartResetPassword(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// getting data from client and binding
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// check email exist
		emailMatchRepo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		emailMatch, err := auxiliary.EmailMatch(emailMatchRepo, input.Email)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
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
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// get userID from DB
		emailRepo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		userID, err := auxiliary.UserIDByEmail(emailRepo, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			if err.Error() == "email not found" {
				c.JSON(http.StatusConflict, gin.H{"error": "email not found"})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
			return
		}

		// add resetPassCode field to User in redis
		if err = UserManageRepo.SavePassCode(resetPassCode, userID); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}
		logs.AuditLogger.Printf("reset pass code sended to user: %s: %s", userID, input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Code sent"})
	}
}
func ApproveResetPassword(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// getting data from client and binding
		var input models.ResetPass
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// check valid password
		passValid, msg := models.ValidPassword(input.NewPassword)
		if !passValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}

		// check email exist in db
		emailMatchRepo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		emailMatch, err := auxiliary.EmailMatch(emailMatchRepo, input.Email)
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}
		if !emailMatch {
			c.JSON(http.StatusConflict, gin.H{"error": "Email not found"})
			return
		}

		// get userID from redis
		emailRepo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		userID, err := auxiliary.UserIDByEmail(emailRepo, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// get reset pass code from redis
		code, err := rdb.HGet(ctx, "user:"+userID, "resetPassCode").Result()
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
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
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
			input.NewPassword = string(hashedPassword)

			// update userVersion
			userVersion := uuid.New().String()

			// save new pass and userVersion in redis
			err = rdb.HSet(ctx, "user:"+userID, map[string]interface{}{
				"password":    input.NewPassword,
				"userVersion": userVersion,
			}).Err()
			if err != nil {
				logs.ErrorLogger.Error(err)
				logs.Error.Println(err)
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			// delete resetPassCode field from redis
			err = rdb.HDel(ctx, "user:"+userID, "resetPassCode").Err()
			if err != nil {
				logs.ErrorLogger.Error(err)
				logs.Error.Println(err)
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
			logs.AuditLogger.Printf("user: %s: %s reset password", userID, input.Email)
			c.JSON(http.StatusOK, gin.H{"msg": "Password updated successfully"})
		}
	}
}
