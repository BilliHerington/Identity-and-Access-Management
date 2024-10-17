package authJWT

import (
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJWT(c *gin.Context, email string) {

	userID, userVersion, err := JwtRepo.GetDataForJWT(email)
	if err != nil {
		if err.Error() == "user does not exist" {
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
		c.JSON(500, gin.H{"error": "please try again later"})
		return
	}

	// save new JWT in DB
	if err = JwtRepo.SetJWT(email, signedToken); err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		c.JSON(500, gin.H{"error": "please try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": signedToken})
}
