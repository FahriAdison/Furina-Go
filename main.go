package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"furina-bot/lib"
	"furina-bot/plugins/general"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	pluginManager *lib.PluginManager
	errorHandler  *lib.ErrorHandler
	sessionManager *lib.SessionManager
	commandParser *lib.CommandParser
)

func main() {
	// Setup error handler dengan recovery
	defer func() {
		if errorHandler != nil {
			errorHandler.RecoverFromPanic("main")
			errorHandler.Close()
		}
	}()

	// Inisialisasi error handler
	var err error
	errorHandler, err = lib.NewErrorHandler("lib/logs")
	if err != nil {
		fmt.Printf("⚠️ Gagal menginisialisasi error handler: %v\n", err)
		// Lanjutkan tanpa error handler
	} else {
		fmt.Println("✅ Error handler berhasil diinisialisasi")
	}

	// Inisialisasi session manager
	sessionManager, err = lib.NewSessionManager("lib/sessions", errorHandler)
	if err != nil {
		if errorHandler != nil {
			errorHandler.LogError(err, "main.sessionManager")
		}
		panic(fmt.Errorf("failed to initialize session manager: %v", err))
	}
	fmt.Println("✅ Session manager berhasil diinisialisasi")



	// Create context
	ctx := context.Background()

	// Get device store dari session manager
	deviceStore, err := sessionManager.GetFirstDevice(ctx)
	if err != nil {
		if errorHandler != nil {
			errorHandler.LogError(err, "main.getFirstDevice")
		}
		panic(err)
	}

	// Setup client logging (reduced verbosity)
	clientLog := waLog.Stdout("Client", "ERROR", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// Inisialisasi command parser
	commandParser = lib.NewCommandParser(lib.DefaultCommandConfig())
	fmt.Println("✅ Command parser berhasil diinisialisasi")

	// Inisialisasi plugin manager
	pluginManager = lib.NewPluginManager(client)
	
	// Daftarkan plugin
	registerPlugins()
	fmt.Println("✅ Plugin manager berhasil diinisialisasi")

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
			if errorHandler != nil {
				errorHandler.LogError(err, "main.connect")
			}
			panic(fmt.Errorf("failed to connect: %v", err))
		}

		fmt.Println("📱 Meminta kode pairing...")
		
		// Request pairing code
		code, err := client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			if errorHandler != nil {
				errorHandler.LogError(err, "main.pairPhone")
			}
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
			if errorHandler != nil {
				errorHandler.LogError(err, "main.reconnect")
			}
			panic(fmt.Errorf("failed to connect: %v", err))
		}
	}

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nMenghentikan bot...")
	
	client.Disconnect()
	fmt.Println("👋 Bot berhasil dihentikan")
}

// registerPlugins mendaftarkan semua plugin yang tersedia
func registerPlugins() {
	var registeredPlugins []string
	
	// Daftarkan plugin ping dari folder general
	pingPlugin := general.NewPingPlugin()
	pluginManager.RegisterPlugin(pingPlugin)
	registeredPlugins = append(registeredPlugins, pingPlugin.GetName())
	
	// Daftarkan plugin help
	helpPlugin := general.NewHelpPlugin()
	pluginManager.RegisterPlugin(helpPlugin)
	registeredPlugins = append(registeredPlugins, helpPlugin.GetName())
	
	// Tampilkan plugins terdaftar dalam satu baris
	fmt.Printf("📦 Plugin terdaftar: %s\n", strings.Join(registeredPlugins, ", "))
}

func eventHandler(evt interface{}) {
	// Tambahkan recovery untuk event handler
	defer func() {
		if errorHandler != nil {
			errorHandler.RecoverFromPanic("eventHandler")
		}
	}()

	switch v := evt.(type) {
	case *events.Message:
		// Handle incoming messages
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
			messageText := v.Message.GetConversation()
			senderJID := v.Info.Sender
			
			fmt.Printf("📨 Pesan dari %s: %s\n", senderJID, messageText)
			
			// Log pesan ke error handler
			if errorHandler != nil {
				errorHandler.LogInfo(fmt.Sprintf("Message from %s: %s", senderJID, messageText), "eventHandler")
			}
			
			// Cek apakah pesan adalah command
			if commandParser.IsCommand(messageText) {
				// Teruskan ke plugin manager
				if err := pluginManager.HandleMessage(v); err != nil {
					fmt.Printf("❌ Error handling command: %v\n", err)
					if errorHandler != nil {
						errorHandler.LogError(err, "eventHandler.pluginManager")
					}
				}
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
		fmt.Println("🎯 Prefix command: ! (contoh: !ping)")
		fmt.Println("⚡ Tekan Ctrl+C untuk menghentikan bot")
		
		if errorHandler != nil {
			errorHandler.LogInfo("Bot connected successfully", "eventHandler")
		}
		
		// Tampilkan info plugin yang tersedia
		plugins := pluginManager.GetAllPlugins()
		fmt.Printf("📦 %d plugin aktif\n", len(plugins))
		
	case *events.Disconnected:
		fmt.Println("❌ Terputus dari WhatsApp")
		if errorHandler != nil {
			errorHandler.LogInfo("Bot disconnected", "eventHandler")
		}
	case *events.LoggedOut:
		fmt.Println("🚪 Logged out dari WhatsApp")
		if errorHandler != nil {
			errorHandler.LogInfo("Bot logged out", "eventHandler")
		}
	}
}