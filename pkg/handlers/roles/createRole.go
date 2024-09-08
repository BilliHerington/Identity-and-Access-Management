package roles

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"encoding/json"
	"errors"
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
	err := initializers.Rdb.HGetAll(c, roleKey).Err()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"role already exists": input.Name})
		return
	}
	//if err == nil {
	//	c.JSON(http.StatusConflict, gin.H{"error": "role already exist"})
	//	return
	//} else {
	//	logs.Error.Println(err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	// Сериализуем Privileges
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
	}, "role:"+input.Name)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"role created successfully": input.Name})
}
