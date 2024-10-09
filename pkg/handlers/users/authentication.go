package users

import (
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type GetPasswordRepository interface {
	GetPassword(userID string) (string, error)
}

func Authenticate(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get data from client and binding with JSON
		var input models.AuthData
		if err := c.ShouldBind(&input); err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check email exist
		emailMatch, err := auxiliary.EmailMatch(&redisAuxiliaryHandlers.RedisEmailRepo{RDB: rdb}, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if !emailMatch {
			c.JSON(http.StatusConflict, gin.H{"error": "Email does not match"})
			return
		}

		// get userID
		userID, err := auxiliary.UserIDByEmail(&redisAuxiliaryHandlers.RedisUserIDByEmailRepo{RDB: rdb}, input.Email)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			if err.Error() == "email not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// get pass
		repo := &redisUsersHandlers.RedisGetPasswordRepo{RDB: rdb}
		pass, err := repo.GetPassword(userID)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Compare the provided password with the hashed password
		err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(input.Password))
		if err != nil {
			logs.Error.Println(err.Error())
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password does not match"})
			return
		}

		// get userVersion
		userVersion, err := auxiliary.UserVersion(&redisAuxiliaryHandlers.RedisUserVersionRepo{RDB: rdb}, userID)
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			if err.Error() == "userVersion not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		jwtHandlers.UpdateJWT(c, userID, userVersion, input.Email, rdb)
		logs.AuditLogger.Printf("User: %s: %s logged in", userID, input.Email)
	}
}
