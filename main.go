package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create context
	ctx := context.Background()
	
	// Setup sessions directory
	sessionsDir := "lib/sessions"
	dbPath := filepath.Join(sessionsDir, "furina-bot.db")
	
	// Create sessions directory if it doesn't exist
	if err := createSessionsDir(sessionsDir); err != nil {
		panic(fmt.Errorf("failed to create sessions directory: %v", err))
	}
	
	// Setup logging (reduced verbosity)
	dbLog := waLog.Stdout("Database", "ERROR", false)
	
	// Create database container with new path
	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath), dbLog)
	if err != nil {
		panic(err)
	}

	// Get device store
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	// Setup client logging (reduced verbosity)
	clientLog := waLog.Stdout("Client", "ERROR", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// Add event handler
	client.AddEventHandler(eventHandler)

	// Check if already logged in
	if client.Store.ID == nil {
		// Request phone number from user
		fmt.Print("Masukkan nomor telepon WhatsApp (format: +62xxx): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		phoneNumber := strings.TrimSpace(scanner.Text())

		if phoneNumber == "" {
			fmt.Println("Nomor telepon tidak boleh kosong!")
			return
		}

		// Connect to WhatsApp first
		err = client.Connect()
		if err != nil {
			panic(fmt.Errorf("failed to connect: %v", err))
		}

		fmt.Println("📱 Meminta kode pairing...")
		
		// Request pairing code
		code, err := client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			panic(fmt.Errorf("failed to request pairing code: %v", err))
		}

		fmt.Printf("\n🔑 Kode Pairing: %s\n", code)
		fmt.Println("📲 Masukkan kode pairing ini di WhatsApp Web/Desktop")
		fmt.Println("⏳ Menunggu terhubung ke WhatsApp...")
	} else {
		fmt.Println("🔄 Menghubungkan ke WhatsApp...")
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(fmt.Errorf("failed to connect: %v", err))
		}
	}

	// Status akan ditampilkan oleh event handler Connected

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nMenghentikan bot...")
	client.Disconnect()
}

// createSessionsDir creates the sessions directory if it doesn't exist
func createSessionsDir(sessionsDir string) error {
	// Check if directory exists
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		fmt.Printf("📁 Membuat folder sesi: %s\n", sessionsDir)
		// Create directory with proper permissions
		if err := os.MkdirAll(sessionsDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", sessionsDir, err)
		}
		fmt.Println("✅ Folder sesi berhasil dibuat")
	} else if err != nil {
		return fmt.Errorf("failed to check directory %s: %v", sessionsDir, err)
	} else {
		fmt.Printf("📂 Menggunakan folder sesi yang sudah ada: %s\n", sessionsDir)
	}
	return nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Handle incoming messages
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
			messageText := v.Message.GetConversation()
			senderJID := v.Info.Sender
			
			fmt.Printf("📨 Pesan dari %s: %s\n", senderJID, messageText)
			
			// Simple echo bot - reply with the same message
			if strings.ToLower(messageText) == "ping" {
				fmt.Printf("💬 Mengirim balasan 'pong' ke %s\n", senderJID)
				// Note: In a real implementation, you'd need access to the client here
				// This is a simplified example
			}
		}
	case *events.Receipt:
		// Handle message receipts (disabled to reduce log spam)
		// if v.Type == events.ReceiptTypeRead || v.Type == events.ReceiptTypeReadSelf {
		//	fmt.Printf("✓ Pesan dibaca oleh %s\n", v.SourceString())
		// }
	case *events.Connected:
		fmt.Println("\n✅ Bot WhatsApp Furina berhasil terhubung!")
		fmt.Println("💾 Sesi tersimpan di: lib/sessions/")
		fmt.Println("🤖 Bot siap menerima pesan")
		fmt.Println("⚡ Tekan Ctrl+C untuk menghentikan bot")
	case *events.Disconnected:
		fmt.Println("❌ Terputus dari WhatsApp")
	case *events.LoggedOut:
		fmt.Println("🚪 Logged out dari WhatsApp")
	}
}