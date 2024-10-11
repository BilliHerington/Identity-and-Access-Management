package users

import (
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementRepository interface {
	GetPassword(email string) (string, error)
	DeleteUserFromDB(email string) error
	GetAllUsersDataFromDB() (map[string]string, error)
	GetUsersListFromDB() ([]string, error)
	StartUserRegistration(userID, email, verificationCode string) error
	GetVerificationCode(userID string) (string, error)
	SaveUser(userID, email, password, name, role, jwt, userVersion string) error
	SavePassCode(userID, passCode string) error
	GerDataForJWT(email string) (userID string, userVersion string, error error)
}

var UserManageRepo UserManagementRepository

func Authenticate(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.AuthData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "incorrect data format, please check your input data"})
			return
		}

		// get saved pass from db
		savedPass, err := UserManageRepo.GetPassword(input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(500, gin.H{"error": "please try again later"})
			return
		}

		// Compare the provided password with the hashed password
		if err = bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(input.Password)); err != nil {
			logs.Error.Println(err.Error())
			logs.ErrorLogger.Error(err.Error())
			c.JSON(400, gin.H{"error": "Password does not match"})
			return
		}

		// get data to update JWT
		userID, userVersion, err := UserManageRepo.GerDataForJWT(input.Email)

		jwtHandlers.UpdateJWT(c, userID, userVersion, input.Email, rdb)
		logs.AuditLogger.Printf("User: %s: %s logged in", userID, input.Email)
	}
}
