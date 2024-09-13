package main

import (
	"IAM/initializers"
	"IAM/pkg/googleAuth"
	"IAM/pkg/handlers/access"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"os"
)

// TODO: защита JWT
// TODO: структурировать доступы

func init() {
	debugMode := os.Getenv("DEBUG_MODE") == "true"
	logs.InitLoggers(debugMode)
	initializers.LoadEnvVariables()

}
func main() {
	initializers.InitRedis()

	// Создание нового маршрутизатора Gin
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
		protected.GET("/get-all-users-data", middlewares.CheckPrivileges("read"), access.GetAllUsersData)
		protected.DELETE("/delete-user", middlewares.CheckPrivileges("delete"), access.DeleteUser)
		//---roles----
		protected.GET("/get-roles", roles.GetRolesList)
		protected.GET("/get-all-roles-data", roles.GetAllRolesData)
		protected.POST("/assign-role", middlewares.CheckPrivileges("edit"), roles.AssignRole)
		protected.POST("/create-role", middlewares.CheckPrivileges("create"), roles.CreateRole)
		protected.DELETE("/delete-role", middlewares.CheckPrivileges("create"), roles.DeleteRole)
		protected.POST("/redact-role", middlewares.CheckPrivileges("create"), roles.RedactRole)
	}
	err := router.Run()
	if err != nil {
		logs.Error.Fatalf("error run server %v", err)
	}
}
