package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func RedactRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.RolesData
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// redact role
		if err := RoleManageRepo.RedactRoleDB(input.RoleName, input.Privileges); err != nil {
			if err.Error() == "role does not exist" {
				c.JSON(400, gin.H{"error": err})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("%s updated successfully. New privileges: %s", input.RoleName, input.Privileges)
		c.JSON(http.StatusOK, gin.H{"role:" + input.RoleName + " updated successfully. New privileges": input.Privileges})
	}
}
