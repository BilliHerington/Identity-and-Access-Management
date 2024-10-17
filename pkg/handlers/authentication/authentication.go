package authentication

import (
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthManagementRepository interface {
	GetPassword(username string) (string, error)
	StartUserRegistration(email, userID, verificationCode string) error
	GetVerificationCode(verificationCode string) (string, error)
	SaveUser(email, password, name, role, jwt, userVersion string) error
	SavePassCode(email, passCode string) error
	SaveNewUserData(email, password, userVersion string) error
}

var AuthManageRepo AuthManagementRepository

func Authenticate() gin.HandlerFunc {
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
		savedPass, err := AuthManageRepo.GetPassword(input.Email)
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

		authJWT.UpdateJWT(c, input.Email)
		logs.AuditLogger.Printf("User: %s logged in", input.Email)
	}
}
