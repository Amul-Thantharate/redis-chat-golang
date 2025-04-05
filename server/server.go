package server

import (
	"bufio"
	"context"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	ctx       = context.Background()
	rdb       = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	users     = make(map[string]net.Conn) // Online users
	userLock  sync.Mutex
	banned    = make(map[string]bool) // Banned users
	muted     = make(map[string]bool) // Muted users
	adminUser = "admin"               // Only this user is admin
)

// Start the server
func StartServer() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("🔥 Chat server started on :8000 🚀")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

// Handle client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter username: /name <yourname>\n"))

	var username string

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "/name ") {
			username = strings.TrimPrefix(text, "/name ")
			break
		}
		conn.Write([]byte("❌ Invalid command. Use /name <yourname>\n"))
	}

	// Check if banned
	if banned[username] {
		conn.Write([]byte("❌ You are banned from this chat.\n"))
		return
	}

	// Register user
	userLock.Lock()
	users[username] = conn
	userLock.Unlock()
	broadcast("✅ " + username + " joined the chat.")

	// Handle messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)

		if muted[username] {
			conn.Write([]byte("❌ You are muted.\n"))
			continue
		}

		switch {
		case message == "/help":
			showHelp(conn)
		case message == "/list_users":
			listUsers(conn)
		case message == "/history":
			showHistory(conn)
		case strings.HasPrefix(message, "/pm "):
			handlePrivateMessage(username, message)
		case strings.HasPrefix(message, "/ban "):
			if username != adminUser {
				conn.Write([]byte("❌ Only the admin can use /ban.\n"))
				continue
			}
			handleBan(username, message)
		case strings.HasPrefix(message, "/unban "):
			if username != adminUser {
				conn.Write([]byte("❌ Only the admin can use /unban.\n"))
				continue
			}
			handleUnban(username, message)
		case strings.HasPrefix(message, "/mute "):
			if username != adminUser {
				conn.Write([]byte("❌ Only the admin can use /mute.\n"))
				continue
			}
			handleMute(username, message)
		case strings.HasPrefix(message, "/unmute "):
			if username != adminUser {
				conn.Write([]byte("❌ Only the admin can use /unmute.\n"))
				continue
			}
			handleUnmute(username, message)
		case strings.HasPrefix(message, "/broadcast "):
			handleBroadcast(username, message)
		case message == "/exit":
			conn.Write([]byte("Goodbye!\n"))
			return
		default:
			broadcast(username + ": " + message)
			saveMessage(username, message)
		}
	}

	// Remove user on disconnect
	userLock.Lock()
	delete(users, username)
	userLock.Unlock()
	broadcast("❌ " + username + " left the chat.")
}

// 📜 Show chat history from Redis
func showHistory(conn net.Conn) {
	messages, err := rdb.LRange(ctx, "chat_history", 0, -1).Result()
	if err != nil {
		conn.Write([]byte("❌ Failed to load history.\n"))
		return
	}
	conn.Write([]byte("📜 Chat History:\n"))
	for _, msg := range messages {
		conn.Write([]byte(msg + "\n"))
	}
}

// 📝 Save messages to Redis
func saveMessage(username, message string) {
	fullMessage := username + ": " + message
	rdb.RPush(ctx, "chat_history", fullMessage)
	rdb.LTrim(ctx, "chat_history", -50, -1) // Keep last 50 messages
}

// 🔒 Ban a user
func handleBan(admin, message string) {
	target := strings.TrimPrefix(message, "/ban ")
	if target == "" || users[target] == nil {
		users[admin].Write([]byte("❌ User not found.\n"))
		return
	}
	banned[target] = true
	users[target].Write([]byte("❌ You have been banned by " + admin + ".\n"))
	users[target].Close()
	broadcast("🚫 " + target + " was banned by " + admin)
}

// 🔓 Unban a user
func handleUnban(admin, message string) {
	target := strings.TrimPrefix(message, "/unban ")
	if target == "" {
		users[admin].Write([]byte("❌ Specify a username to unban.\n"))
		return
	}
	delete(banned, target)
	broadcast("✅ " + target + " was unbanned by " + admin)
}

// 🤫 Mute a user
func handleMute(admin, message string) {
	target := strings.TrimPrefix(message, "/mute ")
	if target == "" || users[target] == nil {
		users[admin].Write([]byte("❌ User not found.\n"))
		return
	}
	muted[target] = true
	users[admin].Write([]byte("🔇 " + target + " has been muted.\n"))
}

// 🔊 Unmute a user
func handleUnmute(admin, message string) {
	target := strings.TrimPrefix(message, "/unmute ")
	if target == "" {
		users[admin].Write([]byte("❌ Specify a username.\n"))
		return
	}
	delete(muted, target)
	users[admin].Write([]byte("🔊 " + target + " has been unmuted.\n"))
}

// 📩 Private message handler
func handlePrivateMessage(sender, message string) {
	parts := strings.SplitN(strings.TrimPrefix(message, "/pm "), " ", 2)
	if len(parts) < 2 {
		users[sender].Write([]byte("❌ Usage: /pm <user> <message>\n"))
		return
	}
	receiver, msg := parts[0], parts[1]
	if users[receiver] == nil {
		users[sender].Write([]byte("❌ User not found.\n"))
		return
	}

	// Send PM to the receiver
	users[receiver].Write([]byte("📩 [PM from " + sender + "]: " + msg + "\n"))

	// Admin sees all private messages
	if users[adminUser] != nil && sender != adminUser {
		users[adminUser].Write([]byte("👀 [PM: " + sender + " -> " + receiver + "] " + msg + "\n"))
	}
}

// 📜 Show online users
func listUsers(conn net.Conn) {
	conn.Write([]byte("👥 Online Users:\n"))
	for user := range users {
		conn.Write([]byte("- " + user + "\n"))
	}
}

// 🛠 Show help menu
func showHelp(conn net.Conn) {
	helpText := `Commands:
/name <username> - Set your username
/list_users - Show online users
/pm <user> <msg> - Send a private message
/history - View chat history
/ban <user> - (Admin only) Ban a user
/unban <user> - (Admin only) Unban a user
/mute <user> - (Admin only) Mute a user
/unmute <user> - (Admin only) Unmute a user
/broadcast <msg> - Send a global message
/exit - Leave the chat`
	conn.Write([]byte(helpText + "\n"))
}

// 📢 Broadcast message to all users
func handleBroadcast(sender, message string) {
	broadcast("📢 [BROADCAST] " + sender + ": " + strings.TrimPrefix(message, "/broadcast "))
}

// 🔊 Send message to all users
func broadcast(message string) {
	userLock.Lock()
	defer userLock.Unlock()
	for _, conn := range users {
		conn.Write([]byte(message + "\n"))
	}
}
