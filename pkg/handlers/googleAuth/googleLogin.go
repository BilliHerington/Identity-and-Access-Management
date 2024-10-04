package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/handlers/gmail"
	"IAM/pkg/handlers/users"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
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

		// set userData
		email := userInfo.EmailAddresses[0].Value
		name := userInfo.Names[0].DisplayName
		var (
			userID      string
			userVersion string
		)

		// check email exist
		emailMatch, err := auxiliary.EmailMatch(&redisAuxiliaryHandlers.RedisEmailRepo{RDB: rdb}, email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// add user to redis if not exist
		if !emailMatch {
			// create userID
			userID = uuid.New().String()[:8]
			userVersion = uuid.New().String()
			err = users.RegistrationRepository.RegisterUser(&redisAuxiliaryHandlers.RedisRegistrationRepo{RDB: rdb}, userID, email, "", name, "user", "", userVersion)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// send email with GmailAPI after login
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

			// if user user already registered get userID and userVersion
			userID, err = auxiliary.UserIDByEmail(&redisAuxiliaryHandlers.RedisUserIDByEmailRepo{RDB: rdb}, email)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
				return
			}
			userVersion, err = auxiliary.UserVersion(&redisAuxiliaryHandlers.RedisUserVersionRepo{RDB: rdb}, userID)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user version"})
				return
			}
		}
		// update JWT
		jwtHandlers.UpdateJWT(c, userID, email, userVersion, rdb)
		logs.AuditLogger.Printf("user:%s: %s is logged in", userID, email)
	}
}
