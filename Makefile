# linear-cli Makefile

.PHONY: build clean test test-verbose test-crud test-crud-verbose test-live-read install lint fmt deps scrape-docs help

# Build variables
BINARY_NAME=linear-cli
GO_FILES=$(shell find . -type f -name '*.go' | grep -v vendor/)
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
# Inject version into cmd.version (overrides default at build time)
LDFLAGS=-ldflags "-X github.com/roboalchemist/linear-cli/cmd.version=$(VERSION)"

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

# Run smoke tests
test:
	@echo "🧪 Running smoke tests..."
	@./smoke_test.sh

# Run smoke tests with verbose output
test-verbose:
	@echo "🧪 Running smoke tests (verbose)..."
	@bash -x ./smoke_test.sh

# Run CRUD integration tests (live API)
test-crud:
	@echo "Running CRUD integration tests (live API)..."
	@go test -v -run TestCRUD -count=1 -timeout 10m .

# Run CRUD integration tests with log file
test-crud-verbose:
	@echo "Running CRUD integration tests (verbose, with log)..."
	@go test -v -run TestCRUD -count=1 -timeout 10m . 2>&1 | tee crud_test.log

# Run the automated LIVE read-only test suite (no write operations)
test-live-read:
	@echo "Running live read-only tests..."
	@./live_read_test.sh --report live_read_report.json

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Scrape Linear's live docs into docs/linear-source-docs (Python 3, stdlib only)
scrape-docs:
	@echo "📚 Scraping Linear docs..."
	@python3 docs/scrape_linear_docs.py

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
	sudo install -m 755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

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

# Run everything: build, format, lint, test, and install
everything: build fmt lint test install
	@echo "✅ Everything complete!"

# Show help
help:
	@echo "📖 Available targets:"
	@echo "  build            - Build the binary"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run smoke tests"
	@echo "  test-verbose     - Run smoke tests with verbose output"
	@echo "  test-crud        - Run CRUD integration tests (live API)"
	@echo "  test-crud-verbose - Run CRUD tests with log file"
	@echo "  test-live-read   - Run automated live read-only tests"
	@echo "  deps             - Install dependencies"
	@echo "  scrape-docs      - Scrape Linear docs into docs/linear-source-docs"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo "  install          - Install binary to system"
	@echo "  dev-install      - Create development symlink"
	@echo "  build-all        - Cross-compile for all platforms"
	@echo "  release          - Prepare release builds"
	@echo "  run              - Build and run the binary"
	@echo "  everything       - Run build, fmt, lint, test, and install"
	@echo "  help             - Show this help"
