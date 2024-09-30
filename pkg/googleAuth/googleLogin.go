package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/gmail"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"net/http"
)

func GoogleLogin(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := initializers.LoadCredentials()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// getting oauth token from HandleOAuthCallback
		token, err := HandleOAuthCallback(c, config)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// getting user info with oauth token
		userInfo, err := GetUserInfo(token, config)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := userInfo.EmailAddresses[0].Value
		name := userInfo.Names[0].DisplayName
		var userID string
		// checking email exist
		match, err := auxiliary.EmailMatch(email, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx := context.Background()
		// adding user to redis if not exist
		if !match {
			// creating userID
			userID = uuid.New().String()[:8]

			err = rdb.Watch(ctx, func(tx *redis.Tx) error {
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

			// adding email key. Example:		(key) email:example@mail (value (userID)) 12345
			err = rdb.Set(ctx, "email:"+email, userID, 0).Err()
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// sending email with GmailAPI after login
			subject := "Welcome to Our Service!"
			body := fmt.Sprintf("Hello %s, welcome to our service. We're excited to have you!", name)
			err = gmail.SendGmail(token, config, email, subject, body)
			if err != nil {
				logs.Error.Println("Error sending email:", err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send welcome email"})
				return
			}
		} else {
			// if user user already registered getting userID
			userID, err = auxiliary.GetUserIDByEmail(ctx, email, rdb)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
				return
			}
		}

		jwtHandlers.UpdateJWT(c, userID, email, rdb)
		logs.AuditLogger.Printf("user:%s: %s is logged in", userID, email)
	}
}
