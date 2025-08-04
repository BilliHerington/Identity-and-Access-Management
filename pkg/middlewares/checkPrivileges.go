package middlewares

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/rolesServices"
	"IAM/pkg/service/usersServices"

	//"IAM/pkg/ginhandlers/rolesHandlers"
	//"IAM/pkg/handlers/users"
	//"IAM/pkg/logs"
	//"IAM/pkg/models"
	//"github.com/gin-gonic/gin"
	//"net/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ManagementMiddlewareRepository interface {
	GetUserRole(email string) (string, error)
	GetRolePrivileges(role string) ([]string, error)
}

func CheckPrivileges(requiredPrivilege string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get userID from header
		email := c.GetString("email")
		if email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		role, err := usersServices.UserManageRepo.GetUserRole(email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		userPrivileges, err := rolesServices.RoleManageRepo.GetRolePrivileges(role)
		//		logs.Info.Println(userPrivileges)
		//		logs.Info.Println(requiredPrivilege)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}

		// compare user privileges with required privileges
		hasPrivileges := false
		for _, privilege := range userPrivileges {
			if privilege == requiredPrivilege {
				hasPrivileges = true
				break
			}
		}
		// if not match
		if !hasPrivileges {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required privileges"})
			c.Abort()
			return
		}
		c.Next()
	}
}
