package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	isDebugMode bool
)

func init() {
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetDebugMode enables or disables debug logging
func SetDebugMode(enabled bool) {
	isDebugMode = enabled
}

// Debug logs a debug message if debug mode is enabled
func Debug(format string, v ...interface{}) {
	if isDebugMode {
		debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(format, v...))
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}
