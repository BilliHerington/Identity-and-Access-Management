package main

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// TODO: логи
// TODO:

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

	// Маршрут для проверки статуса сервера
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := router.Run()
	if err != nil {
		logs.Error.Fatalf("error run server %v", err)
	}
}
