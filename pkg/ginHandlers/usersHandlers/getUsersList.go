package usersHandlers

import (
	"IAM/pkg/models"
	"IAM/pkg/service/usersServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsersList() gin.HandlerFunc {
	return func(c *gin.Context) {

		msg, err := usersServices.GetUsersListService()
		if err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users id`s": msg})
	}
}
