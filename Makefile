.PHONY: build clean run-server test install help docker-build docker-up docker-down docker-logs docker-init docker-clean

# Build both binaries
build:
	@echo "Building vault-server..."
	@go build -o vault-server ./cmd/vault-server
	@echo "Building vault-cli..."
	@go build -o vault-cli ./cmd/vault-cli
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f vault-server vault-cli
	@rm -rf vault-data/
	@echo "Clean complete!"

# Run the vault server
run-server:
	@echo "Starting vault server..."
	@./vault-server

# Run the vault server in dev mode (custom storage path)
run-dev:
	@echo "Starting vault server in dev mode..."
	@./vault-server -storage ./vault-data-dev

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	@cp vault-server $(GOPATH)/bin/
	@cp vault-cli $(GOPATH)/bin/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Docker: Build image
docker-build:
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker build complete!"

# Docker: Start containers
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Vault is running on http://localhost:8200"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Initialize: make docker-init"
	@echo "  2. Or use: ./vault-docker.sh init"

# Docker: Stop containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped!"

# Docker: View logs
docker-logs:
	@docker-compose logs -f vault

# Docker: Initialize vault
docker-init:
	@echo "Initializing Vault..."
	@docker-compose exec vault ./vault-cli init
	@echo ""
	@echo "IMPORTANT: Save the Root Token and Unseal Key above!"

# Docker: Clean everything (including volumes)
docker-clean:
	@echo "Cleaning Docker containers and volumes..."
	@docker-compose down -v
	@echo "Clean complete!"

# Docker: Restart
docker-restart: docker-down docker-up

# Display help
help:
	@echo "Vault Clone - Makefile commands:"
	@echo ""
	@echo "Local builds:"
	@echo "  make build       - Build both vault-server and vault-cli binaries"
	@echo "  make clean       - Remove build artifacts and vault data"
	@echo "  make run-server  - Run the vault server"
	@echo "  make run-dev     - Run the vault server with dev storage path"
	@echo "  make deps        - Install and tidy Go dependencies"
	@echo "  make install     - Install binaries to GOPATH/bin"
	@echo "  make test        - Run tests"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start containers"
	@echo "  make docker-down    - Stop containers"
	@echo "  make docker-logs    - View container logs"
	@echo "  make docker-init    - Initialize vault in container"
	@echo "  make docker-clean   - Remove containers and volumes"
	@echo "  make docker-restart - Restart containers"
	@echo ""
	@echo "  make help        - Display this help message"
