package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	LogFile *os.File
	Logger  *log.Logger
)

func InitLogger() error {
	logDir := "C:\\Users\\User\\.acrapoint"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFileName := fmt.Sprintf("update_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	logFilePath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	LogFile = file
	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func CloseLogger() {
	if LogFile != nil {
		LogFile.Close()
	}
}

func LogInfo(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

func LogError(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}
