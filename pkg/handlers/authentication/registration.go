package authentication

import (
	"IAM/pkg/handlers/emails"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

func StartRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get data from client and binding with JSON
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// code generation
		verificationCode := GenerateVerificationCode()

		// generation ID
		userID := uuid.New().String()[:8]

		// save user in DB
		if err := AuthManageRepo.StartUserRegistration(userID, input.Email, verificationCode); err != nil {
			if err.Error() == "email already registered" {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// email compose
		subject := "Email verification code"
		body := fmt.Sprintf("Your verification code: %s", verificationCode)

		// send email
		if err := emails.SendEmail(subject, body, input.Email); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("user: %s: %s start registration", userID, input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Verification code sent. Check your email address"})
	}
}

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

func ApproveRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {

		// getting data from client and binding
		var input models.RegisterData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// check valid data
		passValid, msg := models.ValidPassword(input.Password)
		if !passValid {
			c.JSON(400, gin.H{"error": msg})
		}
		nameValid, msg := models.ValidName(input.Name)
		if !nameValid {
			c.JSON(400, gin.H{"error": msg})
		}

		// get verification code from DB
		code, err := AuthManageRepo.GetVerificationCode(input.Email)
		logs.Info.Println(code)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// compare codes
		if code != input.VerificationCode {
			c.JSON(400, gin.H{"error": "Invalid verification code"})
			return
		} else {
			// hashing pass
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}
			input.Password = string(hashedPassword)

			// create other data for saving
			userVersion := uuid.New().String()
			role := "user"
			jwt := ""

			// save user in DB
			if err = AuthManageRepo.SaveUser(input.Email, input.Password, input.Name, role, jwt, userVersion); err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			logs.AuditLogger.Printf("user: %s successfully registered", input.Email)
			c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
		}
	}
}
