package main

import (
	"IAM/initializers"
	"IAM/pkg/ginHandlers"
	"IAM/pkg/logs"
	redisDB "IAM/pkg/repository"
	"IAM/pkg/routes"
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

	// init repository
	Rdb, err := redisDB.InitRedis()
	if err != nil {
		logs.AuditLogger.Error(err)
		logs.Error.Fatalf("redis init fail")
	}

	// set dependency
	ginHandlers.InitRepositories(Rdb)
	service.InitRepositoriesForServices(Rdb)

	// set routes
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.RegisterPublicRoutes(router)
	routes.RegisterProtectedRoutes(router)

	logs.Info.Printf("Identity and Access Management launched successfully. Listening on port:%v", os.Getenv("PORT"))
	err = router.Run()
	if err != nil {
		logs.ErrorLogger.Errorf("error running server %v", err)
		logs.Error.Fatalf("error running server %v", err)
	}
}
