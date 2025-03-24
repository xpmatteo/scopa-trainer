# CLAUDE.md - Scopa Trainer Project Guidelines

## Commands
- Build: `make build`
- Run server: `make serve` or `go run cmd/server/main.go`
- Run all tests: `make test` or `go test ./...`
- Run specific test: `go test ./pkg/path/to/package -run TestName -v`
- Run specific package tests: `go test ./pkg/domain -v`
- Clean: `make clean`

## Go Coding Style
- Use MixedCaps (CamelCase) for names, not underscores
- Package names: short, lowercase, single words
- Interfaces: use -er suffix for single-method interfaces
- Error handling: always check errors, use fmt.Errorf for context
- Tests: use table-driven tests with testify assertions
- Follow domain-driven design (pkg/domain, pkg/application, pkg/adapters)

## Http handlers:
- Define small, focused interfaces for each handler's dependencies
- Create handler factory functions that accept these interfaces
- Return http.HandlerFunc from factory functions, not methods on structs

## Project Structure
- `/cmd`: Application entry points
- `/pkg/domain`: Core domain models and logic
- `/pkg/application`: Application services and business logic
- `/pkg/adapters`: External interfaces (HTTP, templates)
- `/templates`: HTML templates
- `/static`: Static assets (images, CSS, JS)

## Testing Guidelines
- Test files alongside implementation
- Use table-driven tests where applicable
- Use testify/assert for assertions
- HTTP handlers: use httptest package
- Avoid logic in test setup: initialize data structures directly
- Afoid logic in test assertions: make a simple assertion for the values we expect

## Quality Assurance Checklist
Before declaring any task complete:
1. Run all tests with `go test ./...` to verify all tests pass
2. Attempt a full build with `make build` or `go build -o scopa-trainer cmd/server/main.go`
3. If the application has UI changes, manually test the affected functionality
4. Check for edge cases in the implemented functionality
5. Ensure code follows project conventions and style guidelines
