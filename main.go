package main

import (
	"fmt"
	"os"

	"chat-server-redis/client"
	"chat-server-redis/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [server|client]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		server.StartServer()
	case "client":
		client.StartClient()
	default:
		fmt.Println("Invalid command. Use 'server' or 'client'.")
		os.Exit(1)
	}
}
