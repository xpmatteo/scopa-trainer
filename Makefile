.PHONY: test build serve clean repomix

# Variables
APP_NAME := scopa-trainer
BUILD_DIR := ./tmp
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
	air

# Clean removes build artifacts
clean:
	rm -rf $(BUILD_DIR) 

# Save repo contents to an AI-friendly single file
# Install repomix with `npm install -g repomix`
repomix:
	repomix -i spec.md
