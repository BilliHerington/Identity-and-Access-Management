package main

import (
	"IAM/initializers"
	"IAM/pkg/googleAuth"
	"IAM/pkg/handlers/access"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	debugMode := os.Getenv("DEBUG_MODE") == "true"
	logs.InitLoggers(debugMode)
	initializers.LoadEnvVariables()
	initializers.InitLogrus()
}
func main() {
	initializers.InitRedis()
	initializers.Lgrs.Info("This is an info log")

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
		logs.Error.Fatalf("error run server %v", err)
	}
}
