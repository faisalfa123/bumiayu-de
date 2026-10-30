# ==============================================================================
# Bumiayu TUI - Makefile
# ==============================================================================

BINARY_NAME=bumiayu
VERSION=0.9.0-beta.5
BUILD_DIR=dist
CMD_DIR=./cmd/bumiayu

.PHONY: all build build-all build-arm64 build-amd64 clean test run lint

all: clean test build

build:
	@echo "==> Building native binary..."
	go build -ldflags="-s -w -X 'main.Version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

build-all: build-arm64 build-amd64 build-armv7
	@echo "==> Multi-arch builds completed in $(BUILD_DIR)/"

# Target khusus ODROID-C2 (Amlogic S905 ARM64)
build-arm64:
	@echo "==> Building for Linux ARM64 (ODROID-C2)..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildArch=arm64'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

# Target Linux Server standar x86_64
build-amd64:
	@echo "==> Building for Linux AMD64 (x86_64 Server)..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildArch=amd64'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)

build-armv7:
	@echo "==> Building for Linux ARMv7..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildArch=armv7'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv7 $(CMD_DIR)

test:
	@echo "==> Running unit tests..."
	go test -v ./...

run:
	@echo "==> Running Bumiayu TUI..."
	go run $(CMD_DIR)

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
