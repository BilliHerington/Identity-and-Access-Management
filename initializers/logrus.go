package initializers

import (
	"github.com/sirupsen/logrus"
	"os"
)

var InfoLogger = logrus.New()
var ErrorLogger = logrus.New()

func InitLogrus() {
	//setting for InfoLogger
	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open info log file: %v", err)
	}
	InfoLogger.SetOutput(infoFile)
	InfoLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	InfoLogger.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("config/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Lgrs.SetOutput(file)
	} else {
		Lgrs.Fatalf("Failed to log to file, using default stderr")
	}
}
