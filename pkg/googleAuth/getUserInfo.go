package googleAuth

import (
	"IAM/pkg/logs"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

// GetUserInfo using token for given info about user from Google People API
func GetUserInfo(token *oauth2.Token, config *oauth2.Config) (*people.Person, error) {
	ctx := context.Background()

	// create HTTP-client with access token
	client := config.Client(context.Background(), token)

	// create client for working with Google People API
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create people client: %v", err)
	}

	// request to API for getting user info
	person, err := srv.People.Get("people/me").PersonFields("names,emailAddresses").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get people: %v", err)
	}

	if len(person.EmailAddresses) == 0 {
		logs.Error.Println("Email not found in user info")
		return nil, fmt.Errorf("email not found")
	}

	return person, nil
}
