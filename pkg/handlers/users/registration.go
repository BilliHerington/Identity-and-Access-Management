package users

import (
	"IAM/pkg/gmail"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}
func StartRegistration(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get data from client and binding with JSON
		var input struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// check email exist
		emailMatch, err := auxiliary.EmailMatch(input.Email, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		logs.Info.Print(body)

		// sending
		err = gmail.SendEmail(subject, body, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{"error": err})
		}

		// generation ID
		userID := uuid.New().String()[:8]

		// save User in Redis
		ctx := context.Background()
		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
					"id":               userID,
					"email":            input.Email,
					"verificationCode": verificationCode,
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

		// add new EmailKey for User
		err = rdb.Set(ctx, "email:"+input.Email, userID, 0).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs.AuditLogger.Printf("user: %s: %s start registration", userID, input.Email)
		c.JSON(http.StatusOK, gin.H{"msg": "Verification code sent"})
	}
}
func ApproveRegistration(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// getting data from client and binding
		var input models.RegisterData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// get id from redis
		userID, err := auxiliary.GetUserIDByEmail(ctx, input.Email, rdb)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// get verification code from redis
		code, err := rdb.HGet(ctx, "user:"+userID, "verificationCode").Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			input.Password = string(hashedPassword)

			// save User in Redis
			err = rdb.Watch(ctx, func(tx *redis.Tx) error {
				_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
					pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
						"id":       userID,
						"email":    input.Email,
						"name":     input.Name,
						"password": input.Password,
						"role":     input.Role,
						"jwt":      input.JWT,
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
			// delete verificationCode field
			err = rdb.HDel(ctx, "user:"+userID, "verificationCode").Err()
			if err != nil {
				logs.Error.Println(err)
				logs.ErrorLogger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			logs.AuditLogger.Printf("user: %s: %s successfully registered", userID, input.Email)
			c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
		}
	}
}
