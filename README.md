<div align="center">

# Furina-Bot

### Furina Bot | Simple WhatsApp Bot Base

<img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTKSN6ccXis9ImmA6Fabma85-gWMPUEvTupDAupIfO4_xT5seNZTSt7BEU&s=10" width="240" height="240" alt="Furina">

</div>

<div align="center">

[![Open Source](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges) [![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)

<img src="https://raw.githubusercontent.com/andreasbm/readme/master/assets/lines/colored.png"/>

</div>

A simple, lightweight WhatsApp bot base built with Go and WhatsMeow library. Perfect for learning and building your own WhatsApp automation.

## 🚧 Development Status

**This project is currently under active development and will continue to be improved as needed.** New features, bug fixes, and enhancements are being added regularly. Feel free to contribute or suggest improvements!

> **Note**: This is a base template designed to be extended with your own features. The core functionality is stable, but additional features are continuously being developed.

## Features

- **🔐 Pairing Code Authentication**: Easy login using pairing code
- **💾 Session Management**: Automatic session storage in `lib/sessions/`
- **🎯 Simple Message Handler**: Clean base for feature development
- **⚡ Lightweight**: Minimal dependencies and fast startup

## Quick Start

### Prerequisites

- Go 1.24 or higher
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone https://github.com/FahriAdison/Furina-Go.git
cd Furina-Go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build
```

4. Run the bot:
```bash
./Furina-Go
```

### First Run

1. Enter your WhatsApp number (format: +62xxx)
2. Enter the pairing code shown into WhatsApp Web/Desktop
3. Bot will connect and ready to receive messages

## Project Structure

```
furina-bot/
├── main.go              # Main source code
├── go.mod               # Module definition
├── go.sum               # Dependencies checksum
├── lib/
│   └── sessions/        # Session storage (auto-created)
│       ├── furina-bot.db
│       └── furina-bot.db-journal
└── README.md           # Documentation
```

## Authentication

### Pairing Code (Default)
```bash
./Furina-Go
```
1. Enter your phone number when prompted
2. Enter the displayed pairing code in WhatsApp Web/Desktop
3. Bot will connect automatically

## Development

### Adding New Features

This bot serves as a simple base ready for development. You can add new features in the `eventHandler()` function in `main.go`.

### Basic Message Handling

The current implementation includes:
- Simple ping response (logs only)
- Message logging with sender information
- Connection status monitoring
- Clean session management

### Extending the Bot

```go
// Example: Add command handling in eventHandler
case *events.Message:
    if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
        messageText := v.Message.GetConversation()
        senderJID := v.Info.Sender
        
        // Add your custom commands here
        switch strings.ToLower(messageText) {
        case "ping":
            // Reply with pong
        case "help":
            // Show help message
        default:
            // Handle other messages
        }
    }
```

## Configuration

The bot uses default configuration. For customization, you can modify:

- Database path in `sessionsDir` variable
- Logging level in `dbLog` and `clientLog`
- User agent in `PairPhone()` parameters

## Built-in Features

### Session Management
- Sessions stored in `lib/sessions/` directory
- Automatic folder creation on first run
- Database contains login info and device keys
- Persistent connection across restarts

### Clean Logging
- Minimal console output
- Error-level logging only
- Connection status updates
- Message reception logs

## Environment Setup

### For Development
```bash
# Install Go dependencies
go mod tidy

# Run in development mode
go run main.go

# Build for production
go build -o Furina-Go
```

### For Production
```bash
# Build optimized binary
go build -ldflags="-s -w" -o Furina-Go

# Run as service
./Furina-Go
```

## Troubleshooting

### Common Issues

1. **"Failed to connect to WhatsApp"**
   - Check internet connection
   - Verify WhatsApp Web is not active on other devices
   - Try clearing session data: `rm -rf lib/sessions`

2. **"Database permission denied"**
   - Ensure `lib/sessions/` directory is writable
   - Check file permissions: `chmod 755 lib/sessions/`

3. **"Invalid phone number format"**
   - Use international format: +62xxx
   - Include country code
   - Remove spaces and special characters

### Debug Mode

For debugging, you can modify logging levels in `main.go`:
```go
// Change ERROR to DEBUG for verbose logging
dbLog := waLog.Stdout("Database", "DEBUG", true)
clientLog := waLog.Stdout("Client", "DEBUG", true)
```

## Security Features

- **Session Encryption**: WhatsApp sessions are encrypted by default
- **Local Storage**: All data stored locally, no external servers
- **Minimal Permissions**: Only requires basic message access
- **Clean Shutdown**: Proper disconnection on exit

## 🗺️ Development Roadmap

### Currently in Development
- [ ] Enhanced command system with prefix support
- [ ] Plugin architecture for modular features
- [ ] Improved error handling and recovery
- [✓] Better session management

### Planned Features
- [ ] Database integration for user data and settings
- [ ] Group management and admin features
- [ ] Media message handling (images, documents, etc.)
- [ ] Webhook support for external integrations
- [ ] Docker containerization
- [ ] Configuration file support
- [ ] Rate limiting and anti-spam features
- [ ] Multi-device session support

### Future Considerations
- [ ] Web dashboard for bot management
- [ ] Analytics and usage statistics
- [ ] Backup and restore functionality
- [ ] Advanced message filtering

## License

This project is licensed under the GPL-2.0 license - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit pull request

## Acknowledgments

- [WhatsMeow](https://github.com/tulir/whatsmeow) - WhatsApp Web API library
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go

## ⚠️ Disclaimer

This bot is created for educational and development purposes. Make sure to comply with WhatsApp's Terms of Service when using this bot.

**Development Notice**: As this project is under continuous development, some features may change or be deprecated in future versions. Always check the latest documentation and changelog before updating.

---

<div align="center">

**Made with ❤️ By Papah-Chan**

</div>
