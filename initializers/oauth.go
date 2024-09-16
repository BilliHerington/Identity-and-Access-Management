package initializers

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"os"
)

func LoadCredentials() (*oauth2.Config, error) {
	file, err := os.Open("config/credentials.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("failed to close credentials file: %v", err)
		}
	}(file)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	scopes := []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/gmail.send", "https://www.googleapis.com/auth/gmail.compose", "https://mail.google.com/"}
	config, err := google.ConfigFromJSON(fileBytes, scopes...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %v", err)
	}
	return config, nil
}
