package lib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// SessionManager mengelola sesi WhatsApp
type SessionManager struct {
	sessionsDir string
	dbPath      string
	container   *sqlstore.Container
	errorHandler *ErrorHandler
}

// NewSessionManager membuat instance baru SessionManager
func NewSessionManager(sessionsDir string, errorHandler *ErrorHandler) (*SessionManager, error) {
	dbPath := filepath.Join(sessionsDir, "furina-bot.db")
	
	sm := &SessionManager{
		sessionsDir:  sessionsDir,
		dbPath:       dbPath,
		errorHandler: errorHandler,
	}

	if err := sm.initializeSessionsDir(); err != nil {
		return nil, err
	}

	if err := sm.initializeDatabase(); err != nil {
		return nil, err
	}

	return sm, nil
}

// initializeSessionsDir membuat direktori sesi jika belum ada
func (sm *SessionManager) initializeSessionsDir() error {
	if _, err := os.Stat(sm.sessionsDir); os.IsNotExist(err) {
		fmt.Printf("📁 Membuat folder sesi: %s\n", sm.sessionsDir)
		if err := os.MkdirAll(sm.sessionsDir, 0755); err != nil {
			return fmt.Errorf("failed to create sessions directory: %v", err)
		}
		fmt.Println("✅ Folder sesi berhasil dibuat")
		
		if sm.errorHandler != nil {
			sm.errorHandler.LogInfo("Sessions directory created", "SessionManager")
		}
	} else if err != nil {
		return fmt.Errorf("failed to check sessions directory: %v", err)
	} else {
		fmt.Printf("📂 Menggunakan folder sesi yang sudah ada: %s\n", sm.sessionsDir)
	}
	return nil
}

// initializeDatabase menginisialisasi database container
func (sm *SessionManager) initializeDatabase() error {
	ctx := context.Background()
	
	// Setup logging dengan level ERROR untuk mengurangi spam
	dbLog := waLog.Stdout("Database", "ERROR", false)
	
	// Buat database container
	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", sm.dbPath), dbLog)
	if err != nil {
		if sm.errorHandler != nil {
			sm.errorHandler.LogError(err, "SessionManager.initializeDatabase")
		}
		return fmt.Errorf("failed to create database container: %v", err)
	}

	sm.container = container
	
	if sm.errorHandler != nil {
		sm.errorHandler.LogInfo("Database container initialized", "SessionManager")
	}
	
	return nil
}

// GetFirstDevice mendapatkan device store pertama
func (sm *SessionManager) GetFirstDevice(ctx context.Context) (*store.Device, error) {
	device, err := sm.container.GetFirstDevice(ctx)
	if err != nil {
		if sm.errorHandler != nil {
			sm.errorHandler.LogError(err, "SessionManager.GetFirstDevice")
		}
		return nil, fmt.Errorf("failed to get first device: %v", err)
	}
	return device, nil
}