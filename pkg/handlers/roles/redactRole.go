package roles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func RedactRole(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.RolesData
		// get data from client and binding with JSON
		if err := c.ShouldBindJSON(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roleKey := "role:" + input.RoleName
		marshalPrivileges, err := json.Marshal(input.Privileges)

		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = rdb.HSet(c, roleKey, "privileges", marshalPrivileges).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs.AuditLogger.Printf("%s updated successfully. New privileges: %s", input.RoleName, input.Privileges)
		c.JSON(http.StatusOK, gin.H{roleKey + " updated successfully. New privileges": input.Privileges})
	}
}
