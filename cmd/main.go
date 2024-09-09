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
	router.GET("/users", access.GetUsers)

	//маршруты для управления ролями
	//router.POST("/create-role", roles.CreateRole)
	router.GET("/roles", roles.GetRoles)
	router.POST("/assign-role", roles.AssignRole)

	router.Use(middlewares.AuthMiddleware())
	//Пример защищенного маршрута, который требует привилегию "create"
	router.POST("/create-role", middlewares.CheckPrivileges("create"), func(c *gin.Context) {}, roles.CreateRole)
	router.DELETE("/delete-user", middlewares.CheckPrivileges("delete"), func(c *gin.Context) {}, access.DeleteUser)
	err := router.Run()
	if err != nil {
		logs.Error.Fatalf("error run server %v", err)
	}
}
