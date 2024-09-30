package logs

import (
	"github.com/fatih/color"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func InitCodeLoggers() {
	color.NoColor = false // disable color auto setting

	infoColor := color.New(color.FgCyan).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()

	// loggers with color messages
	Info = log.New(os.Stdout, infoColor("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, errorColor("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)

}
