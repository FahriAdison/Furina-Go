package lib

import (
	"strings"
)

// CommandConfig konfigurasi untuk sistem command
type CommandConfig struct {
	Prefix      string
	CaseSensitive bool
}

// DefaultCommandConfig konfigurasi default
func DefaultCommandConfig() *CommandConfig {
	return &CommandConfig{
		Prefix:        "!",
		CaseSensitive: false,
	}
}

// CommandParser untuk parsing command dari pesan
type CommandParser struct {
	config *CommandConfig
}

// NewCommandParser membuat instance baru CommandParser
func NewCommandParser(config *CommandConfig) *CommandParser {
	if config == nil {
		config = DefaultCommandConfig()
	}
	return &CommandParser{config: config}
}

// ParseCommand mengurai pesan menjadi command dan arguments
func (cp *CommandParser) ParseCommand(message string) (command string, args []string, isCommand bool) {
	message = strings.TrimSpace(message)
	
	// Cek apakah pesan dimulai dengan prefix
	if !strings.HasPrefix(message, cp.config.Prefix) {
		return "", nil, false
	}

	// Hapus prefix
	message = message[len(cp.config.Prefix):]
	
	// Split menjadi command dan arguments
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "", nil, false
	}

	command = parts[0]
	if !cp.config.CaseSensitive {
		command = strings.ToLower(command)
	}

	args = parts[1:]
	return command, args, true
}

// IsCommand mengecek apakah pesan adalah command
func (cp *CommandParser) IsCommand(message string) bool {
	_, _, isCommand := cp.ParseCommand(message)
	return isCommand
}