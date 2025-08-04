package authenticationServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/service/emailsServices"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"time"
)

func StartUserRegistration(inputUserData models.EmailData) error {
	// code generation
	verificationCode := GenerateVerificationCode()

	// generation ID
	userID := uuid.New().String()[:8]

	// email compose
	subject := "Email verification code"
	body := fmt.Sprintf("Your verification code: %s", verificationCode)

	// send email
	if os.Getenv("USE_TEST_MODE_WITHOUT_GOOGLE") == "true" {
		logs.Info.Printf("You using TEST_MODE_WITHOUT_GOOGLE.\n%v", body)
	} else {
		if err := emailsServices.SendEmail(subject, body, inputUserData.Email); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err)
			return err
		}
	}

	// save user in DB
	if err := AuthManageRepo.StartUserRegistration(userID, inputUserData.Email, verificationCode); err != nil {
		if errors.Is(err, models.ErrEmailAlreadyRegistered) {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}

	logs.AuditLogger.Printf("user: %s: %s start registration", userID, inputUserData.Email)
	return nil
}

func ApproveUserRegistration(inputUserData models.RegisterData) error {
	// get verification code from DB
	code, err := AuthManageRepo.GetVerificationCode(inputUserData.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}

	// compare codes
	if code != inputUserData.VerificationCode {
		err = models.ErrInvalidVerificationCode
		logs.AuditLogger.Printf("user: %s userd invlid verification code", inputUserData.Email)
		return err
	} else {
		// hashing pass
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputUserData.Password), bcrypt.DefaultCost)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			return err
		}
		inputUserData.Password = string(hashedPassword)

		// create other data for saving
		userVersion := uuid.New().String()
		role := "user"
		jwt := ""

		// save user in DB
		if err = AuthManageRepo.SaveUser(inputUserData.Email, inputUserData.Password, inputUserData.Name, role, jwt, userVersion); err != nil {
			if errors.Is(err, models.ErrUserDoesNotExist) {
				return err
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
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
		logs.AuditLogger.Printf("user: %s successfully registered", inputUserData.Email)
	}
	return nil
}

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}
