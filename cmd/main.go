package main

import (
	"IAM/initializers"
	"IAM/pkg/googleAuth"
	"IAM/pkg/handlers/access"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"os"
)

func init() {
	debugMode := os.Getenv("DEBUG_MODE") == "true"
	logs.InitCodeLoggers(debugMode)
	logs.InitFileLoggers()
	initializers.LoadEnvVariables()
}
func SendEmail() error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "kukuruza774@gmail.com")
	mailer.SetHeader("To", "recipient-email@example.com")
	mailer.SetHeader("Subject", "Test Email")
	mailer.SetBody("text/plain", "This is a test email")

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-app-password")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}

	return nil
}
func main() {
	if err := SendEmail(); err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	initializers.InitRedis()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/")
	{
		public.GET("/oauth", googleAuth.OauthRedirect)
		public.GET("auth/callback", googleAuth.GoogleLogin)

		public.POST("/register", access.Registration)
		public.POST("/auth", access.Authenticate)
	}
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware())
		//---users----
		protected.GET("/get-users", access.GetUsersList)
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges(models.Privileges.GetUserData), access.GetAllUsersData)
		protected.DELETE("/delete-user", middlewares.CheckPrivileges(models.Privileges.DeleteUser), access.DeleteUser)
		//---roles----
		protected.GET("/get-roles", roles.GetRolesList)
		protected.GET("/get-all-roles-data", roles.GetAllRolesData)
		protected.POST("/assign-role", middlewares.CheckPrivileges(models.Privileges.CreateRole), roles.AssignRole)
		protected.POST("/create-role", middlewares.CheckPrivileges(models.Privileges.CreateRole), roles.CreateRole)
		protected.DELETE("/delete-role", middlewares.CheckPrivileges(models.Privileges.DeleteRole), roles.DeleteRole)
		protected.POST("/redact-role", middlewares.CheckPrivileges(models.Privileges.CreateRole), roles.RedactRole)
	}
	err := router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}

}
