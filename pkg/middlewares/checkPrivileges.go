package middlewares

import (
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"IAM/pkg/logs"
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

		role, err := users.UserManageRepo.GetUserRole(email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}
		privileges, err := roles.RoleManageRepo.GetRolePrivileges(role)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// compare user privileges with required privileges
		hasPrivileges := false
		for _, privilege := range privileges {
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
