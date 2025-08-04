package usersHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/usersServices"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type UserManagementRepository interface {
//	DeleteUserFromDB(email string) error
//	GetAllUsersDataFromDB() ([]map[string]string, error)
//	GetUsersListFromDB() ([]string, error)
//	GetUserRole(email string) (string, error)
//}

//var UserManageRepo UserManagementRepository

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.EmailData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(400, gin.H{"error": models.ErrIncorrectDataFormat})
			return
		}
		msg, err := usersServices.DeleteUserService(input)
		if err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				c.JSON(400, gin.H{"error": models.ErrUserDoesNotExist})
				return
			}
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": msg})
	}
}
