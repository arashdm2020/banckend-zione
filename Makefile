.PHONY: help build run dev test lint clean swagger docs

# Variables
GO=go
BINARY_NAME=zione-api
MAIN_FILE=cmd/api/main.go
BUILD_DIR=build

# Help
help:
	@echo "Available commands:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application locally"
	@echo "  make dev              - Run the application in development mode with hot reloading"
	@echo "  make test             - Run unit tests"
	@echo "  make lint             - Run linters"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make docs             - Generate API documentation"

# Build the application
build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Run the application locally
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run the application in development mode with hot reloading
dev:
	air -c .air.toml

# Run unit tests
test:
	$(GO) test -v ./internal/...

# Run linters
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -rf tmp

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Generate API documentation
docs: swagger
	@echo "API documentation generated at http://localhost:3000/swagger/index.html"