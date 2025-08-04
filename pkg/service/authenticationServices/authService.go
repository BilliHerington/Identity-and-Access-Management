package authenticationServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthManagementRepository interface {
	GetPassword(username string) (string, error)
	StartUserRegistration(email, userID, verificationCode string) error
	GetVerificationCode(verificationCode string) (string, error)
	SaveUser(email, password, name, role, jwt, userVersion string) error
	SavePassCode(email, passCode string) error
	SaveNewUserData(email, password, userVersion string) error
	DeleteVerificationCode(email string) error
}

var AuthManageRepo AuthManagementRepository

func AuthenticateUser(inputUserData models.AuthData) error {
	// get saved pass from db
	savedPass, err := AuthManageRepo.GetPassword(inputUserData.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}

	// Compare the provided password with the hashed password
	if err = bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(inputUserData.Password)); err != nil {
		return models.ErrPasswordMismatch
	}
	return nil
}
