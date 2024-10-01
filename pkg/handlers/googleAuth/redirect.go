package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

// OauthRedirect redirecting user to URL specified in Credentials.json
func OauthRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// getting config from Credentials.json
		config, err := initializers.LoadCredentials()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, authURL)
	}
}
