package googleAuth

import (
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/models"
	"IAM/pkg/service/authenticationServices/googleAuthService"
	"github.com/gin-gonic/gin"
)

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlCode := c.Query("code")
		userEmail, err := googleAuthService.GoogleLoginUser(urlCode)
		if err != nil {
			c.JSON(500, gin.H{"error": models.ErrInternalServerError})
			return
		}
		// update JWT
		authJWT.UpdateJWT(c, userEmail)
	}
}
