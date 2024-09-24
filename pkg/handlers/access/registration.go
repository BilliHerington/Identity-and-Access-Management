package access

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Registration(c *gin.Context) {
	var input models.RegisterData

	// getting data from client and binding
	if err := c.ShouldBind(&input); err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check email exist
	emailMatch, err := handlers.EmailMatch(input.Email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if emailMatch {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// hashing pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.Password = string(hashedPassword)

	// generation ID
	userID := uuid.New().String()[:8]

	// save User in Redis
	ctx := context.Background()
	err = initializers.Rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
				"id":       userID,
				"email":    input.Email,
				"name":     input.Name,
				"password": input.Password,
				"role":     input.Role,
				"jwt":      input.JWT,
			})
			pipe.SAdd(ctx, "users", userID)
			return nil
		})
		return err
	}, "user:"+userID)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add new EmailKey for User
	err = initializers.Rdb.Set(ctx, "email:"+input.Email, userID, 0).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
