package roles

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RedactRole(c *gin.Context) {
	var input models.RolesData
	if err := c.ShouldBindJSON(&input); err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roleKey := "role:" + input.Name
	marshalPrivileges, err := json.Marshal(input.Privileges)
	logs.Info.Printf("priv:%s", marshalPrivileges)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = initializers.Rdb.HSet(c, roleKey, "privileges", marshalPrivileges).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{roleKey + " updated successfully. New privileges": input.Privileges})
}
