package rolesHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/rolesServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AssignRole set role for user
func AssignRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.UserRoleData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println("error binding data", err)
			logs.ErrorLogger.Error("error binding data", err.Error())
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}

		msg, err := rolesServices.AssignRoleService(input)
		if err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) || errors.Is(err, models.ErrRoleDoesNotExist) {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": msg})
		return
	}
}
