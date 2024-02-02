package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	errWriter *os.File
	logWriter *os.File
}

// NewLogger creates a new Logger with specified error and log file paths.
func NewLogger(errFilePath, logFilePath string) (*Logger, error) {
	errWriter, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logWriter, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		errWriter: errWriter,
		logWriter: logWriter,
	}, nil
}

// Close closes the log files.
func (logger *Logger) Close() {
	_ = logger.errWriter.Close()
	_ = logger.logWriter.Close()
}

// Error logs the error to the error log file
func (logger *Logger) Error(origin string, err error) {
	if err == nil {
		return
	}

	errorMessage := fmt.Sprintf("[%s] ERROR: (%v) %s\n", time.Now().Format(time.RFC3339), origin, err.Error())

	if _, err := logger.errWriter.WriteString(errorMessage); err != nil {
		log.Printf("Failed to write to error log file: %v\n", err)
	}
}

// Error logs the error to the error log file
func (logger *Logger) Log(message string) {
	logMessage := fmt.Sprintf("[%s] MESSAGE: %s\n", time.Now().Format(time.RFC3339), message)

	if _, err := logger.errWriter.WriteString(logMessage); err != nil {
		log.Printf("Failed to write to error log file: %v\n", err)
	}
}
