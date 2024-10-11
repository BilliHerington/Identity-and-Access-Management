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
)

func GoogleLogin(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := initializers.LoadCredentials()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// getting oauth token from HandleOAuthCallback
		token, err := HandleOAuthCallback(c, config)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// getting user info with oauth token
		userInfo, err := GetUserInfo(token, config)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			if err.Error() == "email not found" {
				c.JSON(400, gin.H{"error": "email not found"})
			} else {
				c.JSON(500, gin.H{"error": "please try again later"})
			}
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
		emailMatch, err := auxiliary.EmailMatch(&redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}, email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// add user to redis if not exist
		if !emailMatch {
			// create necessary user data
			userID = uuid.New().String()[:8]
			userVersion = uuid.New().String()
			role := "user"
			pass, jwt := "", ""

			err = users.RegistrationRepository.SaveUser(&redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}, userID, email, pass, name, role, jwt, userVersion)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			// send email with GmailAPI after login
			subject := "Welcome to Our Service!"
			body := fmt.Sprintf("Hello %s, welcome to our service. We're excited to have you!", name)
			err = gmail.SendGmail(token, config, email, subject, body)
			if err != nil {
				logs.Error.Println("Error sending email:", err)
				logs.ErrorLogger.Error("Error sending email:", err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
		} else {

			// if user user already registered get userID and userVersion
			userID, err = auxiliary.UserIDByEmail(&redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}, email)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
			userVersion, err = auxiliary.UserVersion(&redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}, userID)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
		}
		// update JWT
		jwtHandlers.UpdateJWT(c, userID, email, userVersion, rdb)
		logs.AuditLogger.Printf("user:%s: %s is logged in", userID, email)
	}
}
