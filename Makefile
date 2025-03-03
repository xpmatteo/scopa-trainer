.PHONY: test build serve clean

# Variables
APP_NAME := scopa-trainer
BUILD_DIR := ./build
MAIN_FILE := ./cmd/main.go

# Test runs all tests
test:
	go test ./...

# Build compiles the application
build: clean
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Serve runs the application without building a binary
serve:
	go run $(MAIN_FILE)

# Clean removes build artifacts
clean:
	rm -rf $(BUILD_DIR) 