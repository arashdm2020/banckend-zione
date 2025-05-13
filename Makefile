 .PHONY: help build run dev test test-integration lint clean docker-build docker-run docker-dev swagger docs

# Variables
GO=go
BINARY_NAME=zione-api
MAIN_FILE=cmd/api/main.go
BUILD_DIR=build
DOCKER_IMAGE_NAME=zione-api

# Help
help:
	@echo "Available commands:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application locally"
	@echo "  make dev              - Run the application in development mode with hot reloading"
	@echo "  make test             - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make lint             - Run linters"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run application in Docker"
	@echo "  make docker-dev       - Run application in Docker with hot reloading"
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

# Run integration tests
test-integration:
	docker-compose -f docker-compose.test.yml up --build

# Run linters
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -rf tmp

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

# Run application in Docker
docker-run:
	docker-compose up --build

# Run application in Docker with hot reloading
docker-dev:
	docker-compose -f docker-compose.dev.yml up --build

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Generate API documentation
docs: swagger
	@echo "API documentation generated at http://localhost:8080/swagger/index.html"