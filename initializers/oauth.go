package initializers

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
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
	config, err := google.ConfigFromJSON(fileBytes, people.ContactsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %v", err)
	}
	return config, nil
}
