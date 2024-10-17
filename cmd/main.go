package main

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/handlers/authentication"
	"IAM/pkg/handlers/authentication/googleAuth"
	"IAM/pkg/handlers/common"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"IAM/pkg/repository"
	"github.com/gin-gonic/gin"
)

func init() {
	logs.InitCodeLoggers()
	logs.InitFileLoggers()
	initializers.LoadEnvVariables()
}

// TODO use func in initRedis
// TODO check google login with other accounts

func main() {
	Rdb, err := redisDB.InitRedis()
	if err != nil {
		logs.AuditLogger.Error(err)
		logs.Error.Fatalf("redis init fail")
	}
	handlers.InitRepositories(Rdb)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/")
	{
		public.Use(common.RequestLimiter(5, 30))
		public.GET("/oauth", googleAuth.OauthRedirect())
		public.GET("auth/callback", googleAuth.GoogleLogin())

		public.POST("/registration", authentication.StartRegistration())
		public.POST("/approve-registration", authentication.ApproveRegistration())

		public.POST("/forgetPassword", authentication.StartResetPassword())
		public.POST("/updatePassword", authentication.ApproveResetPassword())

		public.POST("/auth", authentication.Authenticate())
	}
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware(), common.RequestLimiter(15, 30))
		//---Users----
		protected.GET("/get-users", users.GetUsersList())
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges(models.AdminPrivileges.GetUserData), users.GetAllUsersData())
		protected.DELETE("/delete-user", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteUser), users.DeleteUser())
		//---Roles----
		protected.GET("/get-roles", roles.GetRolesList())
		protected.GET("/get-all-roles-data", roles.GetAllRolesData())
		protected.POST("/assign-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), roles.AssignRole())
		protected.POST("/create-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), roles.CreateRole())
		protected.DELETE("/delete-role", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteRole), roles.DeleteRole())
		protected.POST("/redact-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), roles.RedactRole())
	}
	logs.Info.Println("Identity and Access Management is starting")
	err = router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}
}
