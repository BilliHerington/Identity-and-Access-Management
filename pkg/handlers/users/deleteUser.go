package users

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserManagementRepository interface {
	DeleteUserFromDB(email string) error
	GetAllUsersDataFromDB() ([]map[string]string, error)
	GetUsersListFromDB() ([]string, error)
	GetUserRole(email string) (string, error)
}

var UserManageRepo UserManagementRepository

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// deleting userdata from DB
		if err := UserManageRepo.DeleteUserFromDB(input.Email); err != nil {
			if err.Error() == "user does not exist" {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Info("user deleted:" + input.Email)
		c.JSON(http.StatusOK, gin.H{"user deleted": input.Email})
	}
}
