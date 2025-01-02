# Project variables
PROJECT_NAME := myst
VERSION := v0.1.0

# Build directory
BUILD_DIR := build

# Go build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

# Default target OS/ARCH is the current system
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: build build-all clean

# Default build for current platform
build:
	@echo "Building $(PROJECT_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)
	@echo "Done! Binary is in $(BUILD_DIR)/$(PROJECT_NAME)"

# Build for all platforms
build-all:
	@echo "Building $(PROJECT_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64
	
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64
	
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)_darwin_arm64
	
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe
	
	@echo "Done! Binaries are in $(BUILD_DIR)/"

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Done!" 