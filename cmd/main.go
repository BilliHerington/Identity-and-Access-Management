package main

import (
	"IAM/initializers"
	"IAM/pkg/googleAuth"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"github.com/go-redis/redis/v8"

	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
)

// TODO set deleter

// TODO divide errors
func init() {
	logs.InitCodeLoggers()
	logs.InitFileLoggers()
	initializers.LoadEnvVariables()
}

var Rdb *redis.Client

func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/")
	{
		public.GET("/oauth", googleAuth.OauthRedirect())
		public.GET("auth/callback", googleAuth.GoogleLogin(Rdb))

		public.POST("/registration", users.StartRegistration(Rdb))
		public.POST("/approve-registration", users.ApproveRegistration(Rdb))

		public.POST("/forgetPassword", users.StartResetPassword(Rdb))
		public.POST("/updatePassword", users.ApproveResetPassword(Rdb))

		public.POST("/auth", users.Authenticate(Rdb))
	}
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware())
		//---users----
		protected.GET("/get-users", users.GetUsersList(Rdb))
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges(models.AdminPrivileges.GetUserData, Rdb), users.GetAllUsersData(Rdb))
		protected.DELETE("/delete-user", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteUser, Rdb), users.DeleteUser(Rdb))
		//---roles----
		protected.GET("/get-roles", roles.GetRolesList(Rdb))
		protected.GET("/get-all-roles-data", roles.GetAllRolesData(Rdb))
		protected.POST("/assign-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole, Rdb), roles.AssignRole(Rdb))
		protected.POST("/create-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole, Rdb), roles.CreateRole(Rdb))
		protected.DELETE("/delete-role", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteRole, Rdb), roles.DeleteRole(Rdb))
		protected.POST("/redact-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole, Rdb), roles.RedactRole(Rdb))
	}
	err := router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}

}
