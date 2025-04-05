package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func StartClient() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go readMessages(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your username:")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	conn.Write([]byte(username + "\n"))

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "exit" {
			break
		}
		conn.Write([]byte(text + "\n"))
	}
}

func readMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(message)
	}
}
