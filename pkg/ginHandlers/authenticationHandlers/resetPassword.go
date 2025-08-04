package authenticationHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/authenticationServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		// getting data from client and binding
		var input models.EmailData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}

		if err := authenticationServices.ResetUserPass(input); err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": "try again later"})
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
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}
		if err := authenticationServices.ApproveResetUserPass(input); err != nil {
			logs.Info.Printf("returned: %s", err)
			if errors.Is(err, models.ErrUserAlreadyExists) || errors.Is(err, models.ErrUserDoesNotExist) || errors.Is(err, models.ErrInvalidVerificationCode) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			logs.Info.Printf("returned2: %s", err)
			c.JSON(500, gin.H{"error": "try again later"})
			return
		}
		// check valid password
		passValid, msg := models.ValidPassword(input.NewPassword)
		if !passValid {
			c.JSON(400, gin.H{"error": msg})
			return
		}

		c.JSON(200, gin.H{"msg": "Password updated successfully"})
	}
}
