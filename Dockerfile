# Use official Golang image as base
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Install dependencies
RUN go mod tidy

# Build the application
RUN go build -o chat-server ./cmd/server/main.go

# Expose server port
EXPOSE 8000

# Start the server
CMD ["/app/chat-server"]
