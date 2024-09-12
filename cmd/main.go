package main

import (
	"IAM/initializers"
	"IAM/pkg/googleAuth"
	"IAM/pkg/handlers/access"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/logs"
	"IAM/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
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

	config, err := initializers.LoadCredentials()
	if err != nil {
		panic(err)
	}
	router.GET("/oauth", func(c *gin.Context) {
		authURL := googleAuth.GetAuthURL(config)
		c.Redirect(http.StatusFound, authURL)
	})
	router.GET("auth/callback", googleAuth.GoogleLogin)

	//маршруты для входа/регистрации
	router.POST("/register", access.Registration)
	router.POST("/auth", access.Authenticate)
	router.GET("/get-users", access.GetUsersList)
	router.GET("/get-all-users-data", access.GetAllUsersData)
	router.DELETE("/delete-user", access.DeleteUser)

	//маршруты для управления ролями
	//router.POST("/create-role", roles.CreateRole)
	router.GET("/get-roles", roles.GetRolesList)
	router.GET("/get-all-roles-data", roles.GetAllRolesData)
	router.POST("/assign-role", roles.AssignRole)
	router.POST("/redact-role", roles.RedactRole)
	//router.POST("/create-role", roles.CreateRole)
	router.DELETE("/delete-role", roles.DeleteRole)

	//router.Use(middlewares.AuthMiddleware())
	////Пример защищенного маршрута, который требует привилегию "create"
	router.POST("/create-role", middlewares.CheckPrivileges("create"), func(c *gin.Context) {}, roles.CreateRole)
	err = router.Run()
	if err != nil {
		logs.Error.Fatalf("error run server %v", err)
	}
}
