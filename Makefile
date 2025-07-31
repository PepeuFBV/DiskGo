# DiskGo Makefile
# Provides various build targets for the DiskGo application

# Variables
BINARY_NAME=diskgo
SOURCE_FILE=src/main.go
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: all build clean test deps help install uninstall

# Default target
all: clean deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=1 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(SOURCE_FILE)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for distribution (static binary)
build-static:
	@echo "Building static binary..."
	CGO_ENABLED=1 $(GOBUILD) \
		$(LDFLAGS) \
		-tags netgo \
		-installsuffix netgo \
		-o $(BINARY_NAME) \
		$(SOURCE_FILE)
	@echo "✅ Static build complete: $(BINARY_NAME)"

# Build for multiple architectures
build-all: build-linux-amd64 build-linux-arm64 build-linux-386

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(SOURCE_FILE)

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(SOURCE_FILE)

build-linux-386:
	@echo "Building for Linux 386..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-386 $(SOURCE_FILE)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies ready"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install system-wide
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installation complete"

# Uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstallation complete"

# Show help
help:
	@echo "DiskGo Build System"
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-static  - Build static binary"
	@echo "  build-all     - Build for all supported architectures"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  test          - Run tests"
	@echo "  install       - Install system-wide"
	@echo "  uninstall     - Remove system installation"
	@echo "  help          - Show this help message"
