package initializers

import (
	"IAM/pkg/logs"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"os"
)

func LoadCredentials() (*oauth2.Config, error) {

	file, err := os.Open("config/credentials.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logs.Error.Printf("failed to close credentials file: %v", closeErr)
			logs.ErrorLogger.Error(closeErr)
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/gmail.send",
		"https://www.googleapis.com/auth/gmail.compose",
		"https://mail.google.com/",
	}

	config, err := google.ConfigFromJSON(fileBytes, scopes...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return config, nil
}
