package googleAuth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// HandleOAuthCallback check authCode from redirected URl
func HandleOAuthCallback(c *gin.Context, config *oauth2.Config) (*oauth2.Token, error) {
	code := c.Query("code") //получение кода из URL
	if code == "" {
		return nil, errors.New("no code in query")
	}
	// Обмен кода на токен
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}
	return token, nil
}
