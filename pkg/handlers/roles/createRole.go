package roles

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

func CreateRole(c *gin.Context) {
	var input models.RolesData
	if err := c.ShouldBindJSON(&input); err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleKey := "role:" + input.Name
	match, err := handlers.RoleMatch(roleKey)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if match {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role already exists"})
		return
	}

	privilegesJSON, err := json.Marshal(input.Privileges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("Error serializing privileges:", err)
		return
	}
	ctx := context.Background()
	err = initializers.Rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, roleKey, map[string]interface{}{
				"name":       input.Name,
				"privileges": privilegesJSON,
			})
			pipe.SAdd(ctx, "roles", input.Name)
			return nil
		})
		return err
	}, roleKey)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"role created successfully": input.Name})
}
