package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/gmail"
	"IAM/pkg/handlers"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"net/http"
)

func GoogleLogin(c *gin.Context) {
	config, err := initializers.LoadCredentials()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := HandleOAuthCallback(c, config)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := GetUserInfo(token, config)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := handlers.EmailMatch(userInfo.EmailAddresses[0].Value)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userID := uuid.New().String()[:8]
	email := userInfo.EmailAddresses[0].Value
	name := userInfo.Names[0].DisplayName

	if !match {
		ctx := context.Background()
		err = initializers.Rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
					"id":       userID,
					"email":    email,
					"name":     name,
					"password": "",
					"role":     "reader",
					"jwt":      "",
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
		err = initializers.Rdb.Set(ctx, "email:"+email, userID, 0).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Отправка email после успешной регистрации
		subject := "Welcome to Our Service!"
		body := fmt.Sprintf("Hello %s, welcome to our service. We're excited to have you!", name)
		err = gmail.SendEmail(token, config, email, subject, body)
		if err != nil {
			logs.Error.Println("Error sending email:", err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send welcome email"})
			return
		}
	}
	jwtHandlers.UpdateJWT(c, userID, email)
}
