# go-static Makefile

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -ldflags "-X github.com/ahoglund/go-static/cmd/go-static/commands.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(LDFLAGS) -o bin/go-static .

# Build for multiple platforms
.PHONY: build-all
build-all:
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/go-static-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/go-static-linux-arm64 .
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/go-static-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/go-static-darwin-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/go-static-windows-amd64.exe .

# Install to GOPATH/bin
.PHONY: install
install:
	go install $(LDFLAGS) .

# Test the application
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	rm -f go-static go-static-root
	rm -rf test-*

# Development build and test
.PHONY: dev
dev: build
	./bin/go-static init example-dev
	./bin/go-static build example-dev
	@echo "Development build complete. Test with: ./bin/go-static serve example-dev"

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Show version
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  install    - Install to GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  dev        - Development build and test"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  version    - Show version info"
	@echo "  help       - Show this help"