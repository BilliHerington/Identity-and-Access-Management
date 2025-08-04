package rolesHandlers

import (
	"IAM/pkg/models"
	"IAM/pkg/service/rolesServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRolesList() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg, err := rolesServices.GetRolesListService()
		if err != nil {
			if errors.Is(err, models.ErrRolesListEmpty) {
				c.JSON(404, err.Error())
				return
			}
			c.JSON(500, models.ErrInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": msg})
	}
}
