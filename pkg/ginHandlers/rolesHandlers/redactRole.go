package rolesHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/rolesServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RedactRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.RolesData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}
		msg, err := rolesServices.RedactRoleService(input)
		if err != nil {
			if errors.Is(err, models.ErrRoleDoesNotExist) {
				c.JSON(400, gin.H{"error": models.ErrRoleDoesNotExist})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": msg})
	}
}
