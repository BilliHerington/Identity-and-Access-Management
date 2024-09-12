package googleAuth

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/jwtHandlers"
	"IAM/pkg/logs"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"net/http"
)

// GetAuthURL generate URL for auth in Google
func GetAuthURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

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

// GetUserInfo using token for given info about user from Google People API
func GetUserInfo(token *oauth2.Token, config *oauth2.Config) (*people.Person, error) {
	client := config.Client(context.Background(), token) // создание HTTP-клиента с токеном доступа
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client)) // создание клиента для работы с Google People API
	if err != nil {
		return nil, fmt.Errorf("unacle to create people client: %v", err)
	}
	// запрос к API для получения информации о пользователе
	person, err := srv.People.Get("people/me").PersonFields("names,emailAddresses").Do()
	if err != nil {
		return nil, fmt.Errorf("unacle to get people: %v", err)
	}
	return person, nil
}
func GoogleLogin(c *gin.Context) {
	config, err := initializers.LoadCredentials()
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := HandleOAuthCallback(c, config)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := GetUserInfo(token, config)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	match, err := handlers.EmailMatch(userInfo.EmailAddresses[0].Value)
	if err != nil {
		logs.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userID := uuid.New().String()[:8]
	email := userInfo.EmailAddresses[0].Value
	name := userInfo.Names[0].DisplayName
	if !match {
		ctx := context.Background()
		err = initializers.Rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
					"id":       userID,
					"email":    email,
					"name":     name,
					"password": "",
					"role":     "reader",
					"jwt":      "",
				})
				pipe.SAdd(ctx, "users", userID)
				return nil
			})
			return err
		}, "user:"+userID)
		if err != nil {
			logs.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = initializers.Rdb.Set(ctx, "email:"+email, userID, 0).Err()
		if err != nil {
			logs.Error.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	jwtHandlers.UpdateJWT(c, userID, email)
}
