.PHONY: help build test test-integration test-all clean install lint fmt vet coverage run

# Variables
BINARY_NAME=audiohoover
GO=go
GOFLAGS=-v
COVERAGE_FILE=coverage.out

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

install: ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS)

test: ## Run unit tests
	@echo "Running unit tests..."
	$(GO) test $(GOFLAGS) -race -cover ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GO) test $(GOFLAGS) -race -tags=integration ./...

test-all: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	$(GO) test $(GOFLAGS) -race -cover ./...
	$(GO) test $(GOFLAGS) -race -tags=integration ./...

coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	$(GO) test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from: https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

clean: ## Remove build artifacts and test cache
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html
	$(GO) clean -testcache

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Development helpers
dev-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

check: fmt vet test ## Run formatting, vetting, and tests
	@echo "All checks passed!"

.DEFAULT_GOAL := help
