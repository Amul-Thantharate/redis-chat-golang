# Go parameters
BUILD_DIR = bin
MAIN_FILE = main.go
SERVER_FILE = ./server/server.go
CLIENT_FILE = ./client/client.go

# Docker parameters
DOCKER_IMAGE = chat-server
DOCKER_CONTAINER = chat-server-container

# Tools
GOLINT = golangci-lint

.PHONY: all build server client run-server run-client test lint fmt clean docker-build docker-run docker-stop deps redis-up redis-down help

## 🏗️ Build both server and client binaries
build: 
	@echo "🔨 Building server and client..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/server $(MAIN_FILE) $(SERVER_FILE)
	@go build -o $(BUILD_DIR)/client $(MAIN_FILE) $(CLIENT_FILE)

## 🔨 Build only the server
server:
	@echo "🔨 Building server..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/server $(MAIN_FILE) $(SERVER_FILE)

## 🔨 Build only the client
client:
	@echo "🔨 Building client..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/client $(MAIN_FILE) $(CLIENT_FILE)

## 🚀 Run the chat server
run-server: server
	@echo "🚀 Starting the server..."
	@go run $(MAIN_FILE) server

## 💬 Run the chat client
run-client: client
	@echo "💬 Running the chat client..."
	@go run $(MAIN_FILE) client

## ✅ Run tests
test: 
	@echo "🧪 Running tests..."
	@go test -v ./...

## 🔎 Lint the code
lint: 
	@echo "🔍 Running lint checks..."
	@$(GOLINT) run ./...

## 🎨 Format the code
fmt: 
	@echo "🎨 Formatting code..."
	@go fmt ./...

## 📥 Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod tidy
	@go mod download

## 🧹 Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BUILD_DIR)

## 🐳 Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

## 🐳 Run Docker container
docker-run: docker-build
	@echo "🚀 Running chat server in Docker..."
	@docker run -p 8000:8000 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

## 🛑 Stop Docker container
docker-stop:
	@echo "🛑 Stopping Docker container..."
	@docker stop $(DOCKER_CONTAINER) || true
	@docker rm $(DOCKER_CONTAINER) || true

## 🚀 Start Redis using Docker Compose
redis-up:
	@echo "🚀 Starting Redis with Docker Compose..."
	@docker-compose up -d

## 🛑 Stop Redis
redis-down:
	@echo "🛑 Stopping Redis..."
	@docker-compose down

## 🔥 Run everything (Build, Lint, Test, Format)
all: fmt lint build test
	@echo "✅ All checks passed!"

## 📜 Show help message
help:
	@echo "📜 Available commands:"
	@echo ""
	@echo "  make build         🔨 Build both server and client"
	@echo "  make server        🔨 Build only the server"
	@echo "  make client        🔨 Build only the client"
	@echo "  make run-server    🚀 Run the chat server"
	@echo "  make run-client    💬 Run the chat client"
	@echo "  make test          ✅ Run all tests"
	@echo "  make lint          🔎 Run lint checks"
	@echo "  make fmt           🎨 Format the code"
	@echo "  make deps          📥 Download dependencies"
	@echo "  make clean         🧹 Clean build artifacts"
	@echo "  make docker-build  🐳 Build the Docker image"
	@echo "  make docker-run    🚀 Run the chat server in Docker"
	@echo "  make docker-stop   🛑 Stop and remove Docker container"
	@echo "  make redis-up      🚀 Start Redis using Docker Compose"
	@echo "  make redis-down    🛑 Stop Redis"
	@echo "  make all           🔥 Run formatting, linting, build, and tests"
	@echo "  make help          📜 Show this help message"
