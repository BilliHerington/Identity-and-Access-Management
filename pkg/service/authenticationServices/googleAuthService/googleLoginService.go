package googleAuthService

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/service/emailsServices"

	//"IAM/pkg/service/emails"
	"github.com/google/uuid"
)

type GoogleLoginRepository interface {
	SaveUserByGoogle(userID, email, password, name, role, jwt, userVersion string) error
	CheckEmailExist(email string) (bool, error)
}

var GoogleLoginRepo GoogleLoginRepository

func GoogleLoginUser(urlCode string) (string, error) {

	// load config from credentials
	config, err := initializers.LoadCredentials()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}

	// get oauth token from HandleOAuthCallback
	token, err := HandleOAuthCallback(urlCode, config)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}

	// get user info with oauth token
	userInfo, err := GetUserInfo(token, config)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}

	email := userInfo.EmailAddresses[0].Value
	name := userInfo.Names[0].DisplayName

	emailExist, err := GoogleLoginRepo.CheckEmailExist(email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
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
			return "", err
		}

		// create Welcome Email
		to := email
		subject := "Welcome Email"
		body := "You successfully registered in our service"

		if err = emailsServices.SendGmail(token, config, to, subject, body); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			return "", err
		}
		return email, nil
	}

	logs.AuditLogger.Printf("user: %s is logged in", email)
	return email, nil
}
