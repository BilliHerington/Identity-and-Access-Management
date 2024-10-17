package emails

import (
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"strings"
)

// SendGmail sends an email using Gmail API
func SendGmail(token *oauth2.Token, config *oauth2.Config, to, subject, body string) error {
	client := config.Client(context.Background(), token)
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create Gmail client: %v", err)
	}

	// Create email message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n" + body)

	// Encode message into Gmail format
	emailMessage := &gmail.Message{
		Raw: encodeWeb64String([]byte(message.String())),
	}

	// Send email
	_, err = srv.Users.Messages.Send("me", emailMessage).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}

// encodeWeb64String is a helper function for encoding email message in base64
func encodeWeb64String(msg []byte) string {
	encoded := base64.URLEncoding.EncodeToString(msg)
	return strings.TrimRight(encoded, "=")
}
