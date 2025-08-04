package main

import (
	"IAM/initializers"
	"IAM/pkg/ginHandlers"
	"IAM/pkg/ginHandlers/authenticationHandlers"
	"IAM/pkg/ginHandlers/authenticationHandlers/googleAuth"
	"IAM/pkg/ginHandlers/requestLimiterHandler"
	"IAM/pkg/ginHandlers/rolesHandlers"
	"IAM/pkg/ginHandlers/usersHandlers"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	redisDB "IAM/pkg/repository"
	"IAM/pkg/service"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	logs.InitCodeLoggers()
	logs.InitFileLoggers()
	initializers.LoadEnvVariables()
}

func main() {

	Rdb, err := redisDB.InitRedis()
	if err != nil {
		logs.AuditLogger.Error(err)
		logs.Error.Fatalf("redis init fail")
	}
	ginHandlers.InitRepositories(Rdb)
	service.InitRepositoriesForServices(Rdb)
	TestModeWthoutGoogle := os.Getenv("USE_TEST_MODE_WITHOUT_GOOGLE")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/")
	{
		public.Use(requestLimiterHandler.RequestLimiter(5, 30))

		if TestModeWthoutGoogle != "true" {
			public.GET("/oauth", googleAuth.OauthRedirect())
			public.GET("auth/callback", googleAuth.GoogleLogin())
		}

		public.POST("/registration", authenticationHandlers.StartRegistration())
		public.POST("/approve-registration", authenticationHandlers.ApproveRegistration())

		public.POST("/forgetPassword", authenticationHandlers.StartResetPassword())
		public.POST("/updatePassword", authenticationHandlers.ApproveResetPassword())

		public.POST("/auth", authenticationHandlers.Authenticate())
	}
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware(), requestLimiterHandler.RequestLimiter(15, 30))
		//---Users----
		protected.GET("/get-users", usersHandlers.GetUsersList())
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges(models.AdminPrivileges.GetUserData), usersHandlers.GetAllUsersData())
		protected.DELETE("/delete-user", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteUser), usersHandlers.DeleteUser())
		//---Roles----
		protected.GET("/get-roles", rolesHandlers.GetRolesList())
		protected.GET("/get-all-roles-data", rolesHandlers.GetAllRolesData())
		protected.POST("/assign-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.AssignRole())
		protected.POST("/create-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.CreateRole())
		protected.DELETE("/delete-role", middlewares.CheckPrivileges(models.AdminPrivileges.DeleteRole), rolesHandlers.DeleteRole())
		protected.POST("/redact-role", middlewares.CheckPrivileges(models.AdminPrivileges.CreateRole), rolesHandlers.RedactRole())
	}

	logs.Info.Printf("Identity and Access Management launched successfully. Listening on port:%v", os.Getenv("PORT"))
	err = router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}
}
