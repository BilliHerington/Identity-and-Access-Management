package logs

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// creating logs for audit and errors
var (
	ErrorLogger = logrus.New()
	AuditLogger = logrus.New()
)

func InitFileLoggers() {
	// setting logger for errors
	errorFile := &lumberjack.Logger{
		Filename:   "config/logs/error.log",
		MaxSize:    10, // mb
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}
	ErrorLogger.SetReportCaller(true)
	ErrorLogger.SetOutput(errorFile)
	ErrorLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05 | 02.01.06", // setting time format
	})
	ErrorLogger.SetLevel(logrus.ErrorLevel) // set level

	// setting logger for audit
	auditFile := &lumberjack.Logger{
		Filename:   "config/logs/audit.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	AuditLogger.SetReportCaller(true)
	AuditLogger.SetOutput(auditFile)
	AuditLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05 | 02.01.06",
	})
	AuditLogger.SetLevel(logrus.InfoLevel)

}
