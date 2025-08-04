package googleAuthService

import (
	"IAM/pkg/logs"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
)

// HandleOAuthCallback check authCode from redirected URl and return Oauth Token
func HandleOAuthCallback(urlCode string, config *oauth2.Config) (*oauth2.Token, error) {
	//code := c.Query("code") // getting code from URL
	if urlCode == "" {
		return nil, errors.New("no code in query")
	}
	// changing Code for token
	token, err := config.Exchange(context.Background(), urlCode)
	if err != nil {
		logs.Info.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, nil
}
