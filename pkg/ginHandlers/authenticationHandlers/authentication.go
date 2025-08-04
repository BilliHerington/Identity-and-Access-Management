package authenticationHandlers

import (
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/authenticationServices"
	"errors"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.AuthData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}
		err := authenticationServices.AuthenticateUser(input)
		if err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, models.ErrPasswordMismatch) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		authJWT.UpdateJWT(c, input.Email)

		logs.AuditLogger.Printf("User: %s logged in", input.Email)
	}
}
