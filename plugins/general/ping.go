package general

import (
	"fmt"
	"runtime"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"furina-bot/lib"
)

// Variabel global untuk menyimpan waktu start bot
var botStartTime = time.Now()

// PingPlugin adalah plugin untuk command ping
type PingPlugin struct{}

// Pastikan PingPlugin mengimplementasikan interface Plugin
var _ lib.Plugin = (*PingPlugin)(nil)

// NewPingPlugin membuat instance baru PingPlugin
func NewPingPlugin() *PingPlugin {
	return &PingPlugin{}
}

// GetName mengembalikan nama plugin
func (p *PingPlugin) GetName() string {
	return "ping"
}

// GetCommands mengembalikan daftar command yang didukung
func (p *PingPlugin) GetCommands() []string {
	return []string{"ping"}
}

// GetDescription mengembalikan deskripsi plugin
func (p *PingPlugin) GetDescription() string {
	return "Plugin sederhana untuk test koneksi bot"
}

// HandleMessage menangani pesan ping
func (p *PingPlugin) HandleMessage(client *whatsmeow.Client, message *events.Message) error {
	messageText := message.Message.GetConversation()
	senderJID := message.Info.Sender

	// Parse command menggunakan command parser
	commandParser := lib.NewCommandParser(lib.DefaultCommandConfig())
	command, _, isCommand := commandParser.ParseCommand(messageText)
	
	if !isCommand {
		return nil
	}

	var responseText string
	startTime := time.Now()
	
	switch command {
	case "ping":
		// Hitung runtime info
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// Hitung berapa lama bot sudah berjalan
		uptime := time.Since(botStartTime)
		
		responseText = fmt.Sprintf(`🏓 *Ping!* Bot Furina-Go aktif dan siap melayani!

📊 *System Information:*
• 🚀 Runtime: %s
• 💾 Memory Usage: %.2f MB
• 🔧 Go Version: %s
• ⚡ Goroutines: %d
• ⏰ Response Time: %v

✨ Bot berjalan dengan lancar!`,
			uptime.Round(time.Second),
			float64(m.Alloc)/1024/1024,
			runtime.Version(),
			runtime.NumGoroutine(),
			time.Since(startTime),
		)
	default:
		return nil
	}

	// Kirim balasan dengan reply menggunakan helper function
	err := lib.SendReplyMessage(client, message, responseText)

	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan ke %s: %v\n", senderJID, err)
		return err
	}

	return nil
}