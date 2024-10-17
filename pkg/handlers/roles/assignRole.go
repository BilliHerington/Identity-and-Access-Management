package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleManagementRepository interface {
	AssignRoleToUser(email, role string) error
	CreateRole(roleName string, privileges []string) error
	DeleteRole(roleName string) error
	GetAllRolesDataFromDB() ([]map[string]string, error)
	GetRolesListFromDB() ([]string, error)
	RedactRoleDB(role string, privileges []string) error
	GetRolePrivileges(role string) ([]string, error)
}

var RoleManageRepo RoleManagementRepository

// AssignRole set role for user
func AssignRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.UserRoleData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println("error binding data", err)
			logs.ErrorLogger.Error("error binding data", err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// assign role
		if err := RoleManageRepo.AssignRoleToUser(input.Email, input.Role); err != nil {
			if err.Error() == "role does not exist" || err.Error() == "user does not exist" {
				c.JSON(400, gin.H{"error": err})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		text := fmt.Sprintf("role: %s, assign successfully for user: %s", input.Role, input.Email)
		logs.AuditLogger.Printf(text)
		c.JSON(http.StatusOK, gin.H{"message": text})
	}
}
