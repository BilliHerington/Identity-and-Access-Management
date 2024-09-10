package main

import (
	"IAM/initializers"
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

	//маршруты для входа/регистрации
	router.POST("/register", access.Registration)
	router.POST("/auth", access.Authenticate)
	router.GET("get-users", access.GetUsersList)
	router.GET("/get-all-users-data", access.GetAllUsersData)

	//маршруты для управления ролями
	//router.POST("/create-role", roles.CreateRole)
	router.GET("/get-roles", roles.GetRolesList)
	router.GET("/get-all-roles-data", roles.GetAllRolesData)
	router.POST("/assign-role", roles.AssignRole)
	router.POST("/redact-role", roles.RedactRole)
	router.POST("/create-role-no-auth", roles.CreateRole)
	router.DELETE("/delete-user", access.DeleteUser)

	router.Use(middlewares.AuthMiddleware())
	//Пример защищенного маршрута, который требует привилегию "create"
	router.POST("/create-role", middlewares.CheckPrivileges("create"), func(c *gin.Context) {}, roles.CreateRole)
	err := router.Run()
	if err != nil {
		logs.Error.Fatalf("error run server %v", err)
	}
}
