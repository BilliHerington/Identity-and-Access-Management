package logs

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

const (
	LevelInfo = iota
	LevelWarning
	LevelError
)

var loglevel = LevelInfo

// LogEvent write event into log-file
func LogEvent(event string) error {
	//open log file
	f, err := os.OpenFile("config/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			Error.Println(err)
		}
	}(f)
	//create new logger
	logger := log.New(f, "", log.LstdFlags)

	//write event into log
	logger.Println(event)
	return nil
}
func LogrusSettings() {
	logger := logrus.New()

	logger.SetOutput(logrus.InfoLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	file, err := os.OpenFile("config/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(file)

}
