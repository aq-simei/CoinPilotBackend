package logger

import (
	"log"
	"os"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

func Init() {
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	warnLogger = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(format string, v ...interface{}) {
	if debugLogger == nil {
		Init()
	}
	debugLogger.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	if infoLogger == nil {
		Init()
	}
	infoLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	if warnLogger == nil {
		Init()
	}
	warnLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	if errorLogger == nil {
		Init()
	}
	errorLogger.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	if errorLogger == nil {
		Init()
	}
	errorLogger.Fatalf(format, v...)
}
