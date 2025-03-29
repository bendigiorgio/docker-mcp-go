SERVER_DIR=cmd/server
AIR_CONFIG=.air.toml

.PHONY: build
build: 
	@echo "Building the server..."
	@go build -o $(SERVER_DIR)/server $(SERVER_DIR)/main.go

.PHONY: run
run: build
	@echo "Starting the server..."
	@$(SERVER_DIR)/server

.PHONY: dev
dev:
	@echo "Starting the server with air..."
	@air -c $(AIR_CONFIG)

.PHONY: dev-sse
dev-sse:
	@echo "Starting the server with air..."
	@air -c $(AIR_CONFIG) -- -t sse

.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

.PHONY: inspector
inspector:
	@echo "Running inspector..."
	@npx @modelcontextprotocol/inspector node build/index.js

.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  build       - Build the server"
	@echo "  run         - Run the server"
	@echo "  dev         - Run the DEV server with air in STDIO mode"
	@echo "  dev-sse     - Run the DEV server with air in SSE mode"
	@echo "  install     - Install dependencies"
	@echo "  inspector   - Run inspector"
	@echo "  help        - Show this help message"