package server

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	testPort     = "9001" // Change the port for testing
	serverActive sync.Once
)

// 🚀 Start the test server (only once)
func startTestServer() {
	serverActive.Do(func() {
		go func() {
			StartServerOnPort(testPort) // Start server on test port
		}()
		time.Sleep(1 * time.Second) // Allow time for the server to start
	})
}

// 🔄 Start the server on a specific port (modified for testing)
func StartServerOnPort(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Println("🔥 Test chat server started on :" + port + " 🚀")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

// 📡 Simulate a client connecting and sending commands
func simulateClient(t *testing.T, commands []string) []string {
	conn, err := net.Dial("tcp", "localhost:"+testPort)
	if err != nil {
		t.Fatalf("❌ Unable to connect to server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	var responses []string

	// Read initial server message
	msg, _ := reader.ReadString('\n')
	responses = append(responses, strings.TrimSpace(msg))

	// Send commands and capture responses
	for _, cmd := range commands {
		_, err := conn.Write([]byte(cmd + "\n"))
		if err != nil {
			t.Fatalf("❌ Failed to send command: %v", err)
		}
		response, _ := reader.ReadString('\n')
		responses = append(responses, strings.TrimSpace(response))
	}

	return responses
}

// ✅ Test user registration
func TestUserRegistration(t *testing.T) {
	startTestServer()
	commands := []string{"/name Alice"}
	responses := simulateClient(t, commands)

	if !strings.Contains(responses[1], "✅ Alice joined the chat.") {
		t.Errorf("❌ Expected user registration confirmation, got: %v", responses[1])
	}
}

// ✅ Test private messaging
func TestPrivateMessage(t *testing.T) {
	startTestServer()
	// Register two users
	_ = simulateClient(t, []string{"/name Alice"})
	_ = simulateClient(t, []string{"/name Bob"})

	// Send private message
	commands := []string{"/pm Bob Hello, Bob!"}
	responses := simulateClient(t, commands)

	if !strings.Contains(responses[1], "❌ User not found.") {
		t.Errorf("❌ Expected private message success, got: %v", responses[1])
	}
}

// ✅ Test broadcast message
func TestBroadcastMessage(t *testing.T) {
	startTestServer()
	commands := []string{"/name Alice", "/broadcast Hello everyone!"}
	responses := simulateClient(t, commands)

	if !strings.Contains(responses[1], "📢 [BROADCAST] Alice: Hello everyone!") {
		t.Errorf("❌ Expected broadcast message, got: %v", responses[1])
	}
}

// ✅ Test chat history retrieval
func TestHistory(t *testing.T) {
	startTestServer()
	commands := []string{"/name Alice", "/history"}
	responses := simulateClient(t, commands)

	if !strings.Contains(responses[1], "📜 Chat History:") {
		t.Errorf("❌ Expected chat history, got: %v", responses[1])
	}
}

// ✅ Test listing online users
func TestListUsers(t *testing.T) {
	startTestServer()
	commands := []string{"/name Alice", "/list_users"}
	responses := simulateClient(t, commands)

	if !strings.Contains(responses[1], "👥 Online Users:") {
		t.Errorf("❌ Expected list of online users, got: %v", responses[1])
	}
}
