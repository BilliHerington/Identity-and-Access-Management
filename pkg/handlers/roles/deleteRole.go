package roles

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteRole(c *gin.Context) {
	//role := c.Param("role")
	rdb := initializers.Rdb
	ctx := initializers.Ctx
	var input models.DeleteRoleData
	if err := c.ShouldBindJSON(&input); err != nil {
		logs.Error.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := input.Name
	err := rdb.SRem(ctx, "roles", role).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = rdb.Del(ctx, "role:"+role).Err()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}
