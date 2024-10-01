package main

import (
	"IAM/initializers"
	redisDB "IAM/initializers/redisSystem"
	"IAM/pkg/handlers/auxiliary"
	googleAuth2 "IAM/pkg/handlers/googleAuth"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"IAM/pkg/models"
	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/")
	{
		public.Use(auxiliary.RequestLimiter(5, 30, Rdb))
		public.GET("/oauth", googleAuth2.OauthRedirect())
		public.GET("auth/callback", googleAuth2.GoogleLogin(Rdb))

		public.POST("/registration", users.StartRegistration(Rdb))
		public.POST("/approve-registration", users.ApproveRegistration(Rdb))

		public.POST("/forgetPassword", users.StartResetPassword(Rdb))
		public.POST("/updatePassword", users.ApproveResetPassword(Rdb))

		public.POST("/auth", users.Authenticate(Rdb))
	}
	protected := router.Group("/")
	{
		protected.Use(middlewares.AuthMiddleware(Rdb), auxiliary.RequestLimiter(10, 60, Rdb))
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
	logs.Info.Println("Identity and Access Management is starting")
	err = router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}

}
