package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.RolesData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// save role
		if err := RoleManageRepo.CreateRole(input.RoleName, input.Privileges); err != nil {
			if err.Error() == "role already exists" {
				c.JSON(400, gin.H{"error": err})
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("role created successfully %s", input.RoleName)
		c.JSON(http.StatusOK, gin.H{"role created successfully": input.RoleName})
	}
}
