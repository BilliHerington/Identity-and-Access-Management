package authJWT

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJWT(c *gin.Context, email string) {
	//c := gin.Context{}
	userID, userVersion, err := JwtRepo.GetDataForJWT(email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return
	}

	// sign token
	signedToken, err := CreateJWT(email, userID, userVersion)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(500, gin.H{"error": models.ErrInternalServerError})
		return
	}

	// save new JWT in DB
	if err = JwtRepo.SetJWT(email, signedToken); err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(500, gin.H{"error": models.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": signedToken})
}
