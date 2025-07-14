# linctl Makefile

.PHONY: build clean test install lint fmt deps help

# Build variables
BINARY_NAME=linctl
GO_FILES=$(shell find . -type f -name '*.go' | grep -v vendor/)
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	go clean

# Run all tests
test:
	@echo "🧪 Running all tests..."
	go test -v ./tests/...

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	go test -v ./tests/unit/...

# Run integration tests (requires LINEAR_TEST_API_KEY)
test-integration:
	@echo "🧪 Running integration tests..."
	@if [ -z "$$LINEAR_TEST_API_KEY" ]; then \
		echo "⚠️  LINEAR_TEST_API_KEY not set. Skipping integration tests."; \
		echo "   Set it with: export LINEAR_TEST_API_KEY=your-test-key"; \
	else \
		go test -v ./tests/integration/...; \
	fi

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Install binary to system
install: build
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Development installation (symlink)
dev-install: build
	@echo "🔗 Creating development symlink..."
	sudo ln -sf $(PWD)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Cross-compile for multiple platforms
build-all:
	@echo "🌍 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Create release directory
release: clean
	@echo "🚀 Preparing release..."
	mkdir -p dist
	$(MAKE) build-all

# Run the binary
run: build
	./$(BINARY_NAME)

# Show help
help:
	@echo "📖 Available targets:"
	@echo "  build            - Build the binary"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires API key)"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  deps             - Install dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo "  install          - Install binary to system"
	@echo "  dev-install      - Create development symlink"
	@echo "  build-all        - Cross-compile for all platforms"
	@echo "  release          - Prepare release builds"
	@echo "  run              - Build and run the binary"
	@echo "  help             - Show this help"