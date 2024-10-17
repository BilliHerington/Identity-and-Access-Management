package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/handlers/emails"
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoogleLoginRepository interface {
	SaveUserByGoogle(userID, email, password, name, role, jwt, userVersion string) error
	CheckEmailExist(email string) (bool, error)
}

var GoogleLoginRepo GoogleLoginRepository

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// load config from credentials
		config, err := initializers.LoadCredentials()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// get oauth token from HandleOAuthCallback
		token, err := HandleOAuthCallback(c, config)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// get user info with oauth token
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
		email := userInfo.EmailAddresses[0].Value
		name := userInfo.Names[0].DisplayName

		emailExist, err := GoogleLoginRepo.CheckEmailExist(email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}
		if !emailExist {
			// set userData
			userID := uuid.New().String()[:8]
			role := "user"
			password, jwtString := "", ""
			userVersion := uuid.New().String()

			// save user in DB
			if err = GoogleLoginRepo.SaveUserByGoogle(userID, email, password, name, role, jwtString, userVersion); err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			// create Welcome Email
			to := email
			subject := "Welcome Email"
			body := "You successfully registered in our service"
			if err = emails.SendGmail(token, config, to, subject, body); err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
		}

		// update JWT
		authJWT.UpdateJWT(c, email)
		logs.AuditLogger.Printf("user: %s is logged in", email)
	}
}
