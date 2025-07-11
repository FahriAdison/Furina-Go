package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ErrorHandler menangani error dengan logging dan recovery
type ErrorHandler struct {
	logFile *os.File
}

// NewErrorHandler membuat instance baru ErrorHandler
func NewErrorHandler(logDir string) (*ErrorHandler, error) {
	// Buat direktori log jika belum ada
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Buat file log dengan timestamp
	logFileName := fmt.Sprintf("furina-bot-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, logFileName)
	
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	return &ErrorHandler{
		logFile: logFile,
	}, nil
}

// LogError mencatat error ke file log
func (eh *ErrorHandler) LogError(err error, context string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] ERROR in %s: %v\n", timestamp, context, err)
	
	// Tulis ke file log
	if eh.logFile != nil {
		eh.logFile.WriteString(logMessage)
	}
	
	// Juga tampilkan di console
	log.Printf("❌ %s: %v", context, err)
}

// LogInfo mencatat informasi ke file log
func (eh *ErrorHandler) LogInfo(message string, context string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] INFO in %s: %s\n", timestamp, context, message)
	
	// Tulis ke file log
	if eh.logFile != nil {
		eh.logFile.WriteString(logMessage)
	}
}

// RecoverFromPanic menangani panic dan mencatatnya
func (eh *ErrorHandler) RecoverFromPanic(context string) {
	if r := recover(); r != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf("[%s] PANIC in %s: %v\n", timestamp, context, r)
		
		// Tulis ke file log
		if eh.logFile != nil {
			eh.logFile.WriteString(logMessage)
		}
		
		// Tampilkan di console
		log.Printf("💥 PANIC in %s: %v", context, r)
		fmt.Println("🔄 Bot akan mencoba melanjutkan operasi...")
	}
}

// Close menutup file log
func (eh *ErrorHandler) Close() error {
	if eh.logFile != nil {
		return eh.logFile.Close()
	}
	return nil
}