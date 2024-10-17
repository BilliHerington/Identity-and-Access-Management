package authentication

import (
	"IAM/pkg/handlers/emails"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func StartResetPassword() gin.HandlerFunc {
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

		//generate code
		resetPassCode := GenerateVerificationCode()

		// email compose
		subject := "Resetting password"
		body := fmt.Sprintf("Code for resetting:%s", resetPassCode)

		// add resetPassCode field to User in redis
		if err := AuthManageRepo.SavePassCode(input.Email, resetPassCode); err != nil {
			if err.Error() == "user does not exist" {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// sending
		if err := emails.SendEmail(subject, body, input.Email); err != nil {
			logs.ErrorLogger.Errorln(err)
			logs.Error.Println(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("reset pass code sended to user: %s", input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Code sent"})
	}
}
func ApproveResetPassword() gin.HandlerFunc {
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
			c.JSON(400, gin.H{"error": msg})
		}

		// get code from DB
		code, err := AuthManageRepo.GetVerificationCode(input.Email)
		if err != nil {
			if err.Error() == "user does not exist" {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			return
		}

		// compare codes
		if code != input.ResetPassCode {
			c.JSON(400, gin.H{"error": "Code Not Match"})
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
			err = AuthManageRepo.SaveNewUserData(input.Email, input.NewPassword, userVersion)
			if err != nil {
				if err.Error() == "user does not exist" {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				logs.ErrorLogger.Error(err)
				logs.Error.Println(err)
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			logs.AuditLogger.Printf("user: %s reset password", input.Email)
			c.JSON(200, gin.H{"msg": "Password updated successfully"})
		}
	}
}
