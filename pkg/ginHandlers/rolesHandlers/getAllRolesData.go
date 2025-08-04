package rolesHandlers

import (
	"IAM/pkg/models"
	"IAM/pkg/service/rolesServices"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetAllRolesData() gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesData, err := rolesServices.GetAllRolesDataService()
		if err != nil {
			if errors.Is(err, models.ErrRolesListEmpty) {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		c.JSON(200, gin.H{"data": rolesData})
	}
}
