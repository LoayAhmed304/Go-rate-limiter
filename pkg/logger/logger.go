// Package logger provides a simple logging utility with thread safety.
// It supports logging informational and error messages to stdout and stderr respectively.
//
// Used as a singleton instance to ensure consistent logging behavior and format across the application.
// Logger must be initialized before use, typically at the start of the application.
package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	mu    sync.Mutex
	info  log.Logger
	error log.Logger
}

var loggerInstance *Logger = &Logger{}

// Init initializes the logger instance.
func Init() {
	loggerInstance = &Logger{
		mu:    sync.Mutex{},
		info:  *log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime),
		error: *log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	loggerInstance.info.Println("Logger initialized")
}

// LogInfo logs an informational message to stdout.
// With a format of "[INFO]: <date> <time> <message>".
func LogInfo(msg string) {
	if loggerInstance == nil {
		log.Println("Logger is not initialized.")
		return
	}

	loggerInstance.mu.Lock()
	defer loggerInstance.mu.Unlock()

	loggerInstance.info.Println(msg)
}

// LogError logs an error message to stderr.
// With a format of "[ERROR]: <date> <time> <file:line> message".
func LogError(msg string) {
	if loggerInstance == nil {
		log.Println("Logger is not initialized.")
		return
	}

	loggerInstance.mu.Lock()
	defer loggerInstance.mu.Unlock()

	loggerInstance.error.Println(msg)
}
