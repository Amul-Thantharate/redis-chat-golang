# 💬 Redis Chat Server in Go

A real-time chat server implementation using Go and Redis, featuring:
- Multi-client support with concurrent connections
- Redis for message persistence and real-time data
- Command-based chat interface
- Support for both private and broadcast messages
- Admin controls and moderation features
- Scalable architecture with Docker support

## 🌟 Features

- 🔒 User authentication and management
- 💭 Real-time messaging
- 📨 Private messaging support
- 📢 Broadcast messages
- 📜 Chat history with Redis persistence
- 🔨 Moderation tools (ban/mute users)
- 🐳 Docker and Docker Compose support
- 🚀 Easy to deploy and scale

## Documentation
For detailed feature list, please check [Features](./FEATUERS.md)

## 🛠️ Prerequisites

- Go 1.21 or higher
- Redis server
- Docker (optional)
- Make (optional)

## 🔧 Installation

1. Clone the repository:
```bash
git clone https://github.com/Amul-Thantharate/redis-chat-golang.git
cd redis-chat-golang
```

2. Install dependencies:
```bash
make deps
```

3. Start Redis (choose one method):
```bash
# Using Docker Compose
make redis-up

# Or using your local Redis installation
redis-server
```

## 🚀 Running the Application

### Using Make commands:

1. Start the server:
```bash
make run-server
```

2. Start the client (in a new terminal):
```bash
make run-client
```

### Using Go directly:

1. Start the server:
```bash
go run main.go server
```

2. Start the client:
```bash
go run main.go client
```

### Using Docker:

```bash
# Build and run the server
make docker-run

# Stop the server
make docker-stop
```

## 💡 Available Commands

Once connected, users can use the following commands:

| Command | Description |
|---------|-------------|
| `/name <username>` | 👤 Set your username |
| `/pm <user> <message>` | 📩 Send a private message |
| `/broadcast <message>` | 📢 Send a message to all users |
| `/list_users` | 👥 Show online users |
| `/history` | 📜 View chat history |
| `/ban <user>` | 🚫 Ban a user (admin only) |
| `/unban <user>` | ✅ Unban a user (admin only) |
| `/mute <user>` | 🔇 Mute a user (admin only) |
| `/unmute <user>` | 🔊 Unmute a user (admin only) |
| `/help` | ℹ️ Show help menu |
| `/exit` | 👋 Leave the chat |

## 🏗️ Project Structure

```
redis-chat-golang/
├── server/           # Server implementation
├── client/           # Client implementation
├── bin/             # Compiled binaries
├── Dockerfile       # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── Makefile         # Build and run commands
└── README.md        # Documentation
```

## 🧪 Running Tests

```bash
make test
```

## 🔍 Code Quality

```bash
# Run linting
make lint

# Format code
make fmt
```

## 📝 Development Commands

| Command | Description |
|---------|-------------|
| `make build` | 🔨 Build both server and client |
| `make clean` | 🧹 Clean build artifacts |
| `make deps` | 📦 Download dependencies |
| `make all` | 🔥 Run formatting, linting, build, and tests |
| `make help` | 📜 Show all available commands |

## 🐳 Docker Support

The application includes Docker support for easy deployment:

```bash
# Build Docker image
make docker-build

# Run server in Docker
make docker-run

# Stop Docker container
make docker-stop
```

## 🔒 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ⭐️ Show your support

Give a ⭐️ if this project helped you!

## 📞 Contact

- GitHub: [@Amul-Thantharate](https://github.com/Amul-Thantharate)
- Email: amulthantharate@gmail.com

## 🙏 Acknowledgments

- Go Redis Client
- Docker Community
- Open Source Community
