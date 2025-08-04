package authenticationServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/emailsServices"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func ResetUserPass(inputUserData models.EmailData) error {
	//generate code
	resetPassCode := GenerateVerificationCode()

	// email compose
	subject := "Resetting password"
	body := fmt.Sprintf("Code for resetting:%s", resetPassCode)

	// sending mail
	if os.Getenv("USE_TEST_MODE_WITHOUT_GOOGLE") == "true" {
		logs.Info.Printf("You using TEST_MODE_WITHOUT_GOOGLE.\n%v", body)
	} else {
		if err := emailsServices.SendEmail(subject, body, inputUserData.Email); err != nil {
			logs.ErrorLogger.Errorln(err)
			logs.Error.Println(err)
			return err
		}
	}

	// add resetPassCode field to User in repository
	if err := AuthManageRepo.SavePassCode(inputUserData.Email, resetPassCode); err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return err
	}

	return nil
}
func ApproveResetUserPass(inputUserData models.ResetPass) error {
	//logs.Info.Printf("Resetting password for user %s", inputUserData.Email)
	// get code from DB
	code, err := AuthManageRepo.GetVerificationCode(inputUserData.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}

	// compare codes
	if code != inputUserData.ResetPassCode {
		err = models.ErrInvalidVerificationCode
		logs.AuditLogger.Printf("user: %s used invalid verification code", inputUserData.Email)
		return err
	} else {
		// hashing pass
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputUserData.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			return err
		}
		inputUserData.NewPassword = string(hashedPassword)

		// update userVersion
		userVersion := uuid.New().String()

		// save new pass and userVersion in redis
		if err = AuthManageRepo.SaveNewUserData(inputUserData.Email, inputUserData.NewPassword, userVersion); err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				return err
			}
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
			return err
		}
		if err := AuthManageRepo.DeleteVerificationCode(inputUserData.Email); err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				return err
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			return err
		}
		logs.AuditLogger.Printf("user: %s reset password", inputUserData.Email)
	}
	return nil
}
