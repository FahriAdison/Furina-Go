package lib

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

// Plugin interface untuk semua plugin
type Plugin interface {
	// GetName mengembalikan nama plugin
	GetName() string
	
	// GetCommands mengembalikan daftar command yang didukung plugin
	GetCommands() []string
	
	// HandleMessage menangani pesan yang masuk
	HandleMessage(client *whatsmeow.Client, message *events.Message) error
	
	// GetDescription mengembalikan deskripsi plugin
	GetDescription() string
}

// PluginManager mengelola semua plugin
type PluginManager struct {
	plugins       map[string]Plugin
	client        *whatsmeow.Client
	commandParser *CommandParser
}

// NewPluginManager membuat instance baru PluginManager
func NewPluginManager(client *whatsmeow.Client) *PluginManager {
	return &PluginManager{
		plugins:       make(map[string]Plugin),
		client:        client,
		commandParser: NewCommandParser(DefaultCommandConfig()),
	}
}

// RegisterPlugin mendaftarkan plugin baru
func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
	pm.plugins[plugin.GetName()] = plugin
}

// HandleMessage menangani pesan dan meneruskan ke plugin yang sesuai
func (pm *PluginManager) HandleMessage(message *events.Message) error {
	if message.Info.IsFromMe || message.Message.GetConversation() == "" {
		return nil
	}

	messageText := message.Message.GetConversation()
	
	// Parse command menggunakan command parser
	command, _, isCommand := pm.commandParser.ParseCommand(messageText)
	if !isCommand {
		return nil
	}

	// Cari plugin yang menangani command ini
	for _, plugin := range pm.plugins {
		for _, cmd := range plugin.GetCommands() {
			if cmd == command {
				return plugin.HandleMessage(pm.client, message)
			}
		}
	}

	return nil
}

// GetAllPlugins mengembalikan semua plugin yang terdaftar
func (pm *PluginManager) GetAllPlugins() map[string]Plugin {
	return pm.plugins
}

// SendReply mengirim pesan balasan dengan quote/reply
func (pm *PluginManager) SendReply(message *events.Message, responseText string) error {
	senderJID := message.Info.Sender
	chatJID := message.Info.Chat
	messageText := message.Message.GetConversation()
	senderJIDString := senderJID.String()

	// Buat pesan dengan quote/reply
	replyMessage := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &responseText,
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    &message.Info.ID,
				Participant: &senderJIDString,
				QuotedMessage: &waE2E.Message{
					Conversation: &messageText,
				},
			},
		},
	}

	// Kirim balasan dengan reply
	_, err := pm.client.SendMessage(context.Background(), chatJID, replyMessage)
	return err
}

// SendSimpleReply mengirim pesan balasan sederhana dengan quote
func SendReplyMessage(client *whatsmeow.Client, message *events.Message, responseText string) error {
	senderJID := message.Info.Sender
	chatJID := message.Info.Chat
	messageText := message.Message.GetConversation()
	senderJIDString := senderJID.String()

	// Buat pesan dengan quote/reply
	replyMessage := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &responseText,
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    &message.Info.ID,
				Participant: &senderJIDString,
				QuotedMessage: &waE2E.Message{
					Conversation: &messageText,
				},
			},
		},
	}

	// Kirim balasan dengan reply
	_, err := client.SendMessage(context.Background(), chatJID, replyMessage)
	return err
}