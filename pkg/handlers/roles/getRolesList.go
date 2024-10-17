package roles

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRolesList() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := RoleManageRepo.GetRolesListFromDB()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			if err.Error() == "redisRoles not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"redisRoles list": roles})
	}
}
