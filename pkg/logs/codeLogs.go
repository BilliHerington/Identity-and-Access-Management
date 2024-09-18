package logs

import (
	"github.com/fatih/color"
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger // Новый логгер для отладочных сообщений
)

func InitLoggers(debugMode bool) {
	color.NoColor = false // Отключаем автоматическое определение поддержки цвета

	infoColor := color.New(color.FgCyan).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()
	debugColor := color.New(color.BgBlue).SprintFunc()

	// Логгеры с цветными сообщениями
	Info = log.New(os.Stdout, infoColor("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, errorColor("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)

	// Проверяем, включен ли режим отладки

	if debugMode {
		Debug = log.New(os.Stdout, debugColor("DEBUG: "), log.Ldate|log.Ltime|log.Lshortfile)
		Debug.Println("DEBUG MODE ON")
	} else {
		Debug = log.New(io.Discard, "", 0) // Пустой вывод
	}

}
