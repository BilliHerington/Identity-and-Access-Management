package authenticationHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/authenticationServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get data from client and binding with JSON
		var input models.EmailData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}

		err := authenticationServices.StartUserRegistration(input)
		if err != nil {
			if errors.Is(err, models.ErrEmailAlreadyRegistered) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": "try again later"})
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Verification code sent. Check your email address"})
	}
}

func ApproveRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {

		// getting data from client and binding
		var input models.RegisterData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}

		// check valid data
		passValid, msg := models.ValidPassword(input.Password)
		if !passValid {
			c.JSON(400, gin.H{"error": msg})
			return
		}
		nameValid, msg := models.ValidName(input.Name)
		if !nameValid {
			c.JSON(400, gin.H{"error": msg})
			return
		}

		err := authenticationServices.ApproveUserRegistration(input)
		if err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, models.ErrInvalidVerificationCode) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": "try again later"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})

	}
}
