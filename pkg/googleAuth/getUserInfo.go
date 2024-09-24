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
	client := config.Client(context.Background(), token) // создание HTTP-клиента с токеном доступа
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client)) // создание клиента для работы с Google People API
	if err != nil {
		return nil, fmt.Errorf("unable to create people client: %v", err)
	}

	// запрос к API для получения информации о пользователе
	person, err := srv.People.Get("people/me").PersonFields("names,emailAddresses").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get people: %v", err)
	}

	// Логирование информации об электронной почте
	if len(person.EmailAddresses) == 0 {
		logs.Error.Println("Email not found in user info")
		return nil, fmt.Errorf("email not found")
	}

	//logs.Info.Println("User email:", person.EmailAddresses[0].Value)
	return person, nil
}
