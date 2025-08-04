package emailsServices

import (
	"IAM/pkg/logs"
	"net/smtp"
	"os"
)

func SendEmail(subject, body, toEmail string) error {
	// setting SMTP server
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	fromEmail := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")

	// message
	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body)
	// redisAuthentication  SMTP
	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, msg)
	if err != nil {
		logs.Error.Printf("email data:\nhost:port - %s:%s\nauth:%s\nfromEmail:%s\ntoEmail:%s\nmsg:%s\n", smtpHost, smtpPort, auth, fromEmail, toEmail, body)
		return err
	}
	return nil
}
