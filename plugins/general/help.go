package general

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"furina-bot/lib"
)

// HelpPlugin adalah plugin untuk command help
type HelpPlugin struct{}

// Pastikan HelpPlugin mengimplementasikan interface Plugin
var _ lib.Plugin = (*HelpPlugin)(nil)

// NewHelpPlugin membuat instance baru HelpPlugin
func NewHelpPlugin() *HelpPlugin {
	return &HelpPlugin{}
}

// GetName mengembalikan nama plugin
func (h *HelpPlugin) GetName() string {
	return "help"
}

// GetCommands mengembalikan daftar command yang didukung
func (h *HelpPlugin) GetCommands() []string {
	return []string{"menu"}
}

// GetDescription mengembalikan deskripsi plugin
func (h *HelpPlugin) GetDescription() string {
	return "Plugin untuk menampilkan bantuan dan daftar command yang tersedia"
}

// HandleMessage menangani pesan help
func (h *HelpPlugin) HandleMessage(client *whatsmeow.Client, message *events.Message) error {
	messageText := message.Message.GetConversation()
	senderJID := message.Info.Sender

	// Parse command menggunakan command parser
	commandParser := lib.NewCommandParser(lib.DefaultCommandConfig())
	command, _, isCommand := commandParser.ParseCommand(messageText)
	
	if !isCommand {
		return nil
	}

	var responseText string
	
	switch command {
	case "menu":
		responseText = h.generateHelpText()
	default:
		return nil
	}

	// Kirim balasan dengan reply menggunakan helper function
	err := lib.SendReplyMessage(client, message, responseText)

	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan help ke %s: %v\n", senderJID, err)
		return err
	}

	return nil
}

// generateHelpText menghasilkan teks bantuan
func (h *HelpPlugin) generateHelpText() string {
	var help strings.Builder
	
	help.WriteString("🤖 *Furina-Go Bot - Bantuan*\n\n")
	help.WriteString("📋 *Daftar Command yang Tersedia:*\n\n")
	
	// General Commands
	help.WriteString("🔧 *General Commands:*\n")
	help.WriteString("• `!ping` - Cek status bot dan info sistem\n")
	help.WriteString("• `!menu` - Tampilkan bantuan ini\n\n")
	
	// Bot Info
	help.WriteString("ℹ️ *Informasi Bot:*\n")
	help.WriteString("• Prefix: `!` (tanda seru)\n\n")
	
	// Footer
	help.WriteString("✨ *Furina-Go Bot v1.0*\n")
	help.WriteString("🔗 Powered by Papah-Chan\n")
	
	return help.String()
}