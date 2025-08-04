package usersHandlers

import (
	"IAM/pkg/models"
	"IAM/pkg/service/usersServices"
	"github.com/gin-gonic/gin"
)

func GetAllUsersData() gin.HandlerFunc {
	return func(c *gin.Context) {

		allUsersData, err := usersServices.GetAllUsersDataService()
		if err != nil {
			c.JSON(505, gin.H{"error": models.ErrInternalServerError})
		}
		c.JSON(200, allUsersData)
	}
}
