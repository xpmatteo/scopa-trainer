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
bundle:
	repomix -i spec.md

coverage:
	go test -coverprofile /tmp/cover.out ./...
	go tool cover -html=cover.out

# Incomplete!!!
save-prompts:
	cat ~/.claude.json | jq '.projects."/Users/matteo/dojo/2025-03-01-scopa-trainer-take-2/scopa"' > human-docs/claude.json
	pushd /tmp; \
		cursor-export -w '/Users/matteo/Library/Application Support/Cursor/User/workspaceStorage' \
		popd \
		cp -r /tmp/cursor-export-output/markdown/scopa human-docs
