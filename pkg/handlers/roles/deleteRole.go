package roles

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input struct {
			RoleName string `json:"role_name"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Print(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// delete role
		if err := RoleManageRepo.DeleteRole(input.RoleName); err != nil {
			if err.Error() == "role does not exist" {
				c.JSON(400, gin.H{"error": err})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("role deleted successfully: %s", input.RoleName)
		c.JSON(http.StatusOK, gin.H{"role deleted successfully: ": input.RoleName})
	}
}
