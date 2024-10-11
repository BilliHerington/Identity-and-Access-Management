package users

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/handlers/gmail"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

func StartRegistration(rdb *redis.Client) gin.HandlerFunc {
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

		// check email exist
		repo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		emailMatch, err := auxiliary.EmailMatch(repo, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		} else if emailMatch {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		// code generation
		verificationCode := GenerateVerificationCode()

		// email compose
		subject := "Email verification code"
		body := fmt.Sprintf("Your verification code:%s", verificationCode)

		// sending
		if err = gmail.SendEmail(subject, body, input.Email); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err)
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// generation ID
		userID := uuid.New().String()[:8]

		// save user in DB
		if err = UserManageRepo.StartUserRegistration(userID, input.Email, verificationCode); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		logs.AuditLogger.Printf("user: %s: %s start registration", userID, input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Verification code sent"})
	}
}

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

func ApproveRegistration(rdb *redis.Client) gin.HandlerFunc {
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

		// get userID from DB
		userIDRepo := &redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}
		userID, err := auxiliary.UserIDByEmail(userIDRepo, input.Email)
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

		// get verification code from DB
		code, err := UserManageRepo.GetVerificationCode(userID)
		logs.Info.Println(code)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// compare codes
		if code != input.VerificationCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid verification code"})
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
			if err = UserManageRepo.SaveUser(userID, input.Email, input.Password, input.Name, role, jwt, userVersion); err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(500, gin.H{"error": "please try again later"})
				return
			}

			logs.AuditLogger.Printf("user: %s: %s successfully registered", userID, input.Email)
			c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
		}
	}
}
