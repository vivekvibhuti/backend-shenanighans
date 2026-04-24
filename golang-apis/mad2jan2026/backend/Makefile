.PHONY: help build run clean test lint fmt vet tidy deps migrate docker-build

# Binary name
BINARY_NAME=api
BINARY_DIR=bin
MAIN_PATH=./cmd/api

# Go commands
GO=go
GOCACHE=$(shell $(GO) env GOCACHE)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: build ## Build and run the binary
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_DIR)/$(BINARY_NAME)

clean: ## Remove build artifacts and cache
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -rf $(GOCACHE)
	@$(GO) clean -cache -testcache

fmt: ## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running vet..."
	@$(GO) vet ./...

tidy: ## Tidy go.mod
	@echo "Tidying dependencies..."
	@$(GO) mod tidy

deps: tidy ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GO) mod download

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v -race -cover ./...

test-cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./...

lint-fix: ## Run linter with auto-fix
	@echo "Running linter with fixes..."
	@golangci-lint run --fix ./...

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run scripts/migrate.go

docker-build: ## Build Docker image (not using docker images right now, will not use docker builds on this laptop:/)
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

build-linux: ## Build for Linux (cross-compile)
	@echo "Cross-compiling for Linux..."
	@mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
