# Workspace: /Users/matteo/dojo/2025-03-01-scopa-trainer-take-2/scopa

Last Modified: 2025-03-24 16:55:39

## Go Best Practices MDC Rules Creation

**User**:

@Effective Go  @000-cursor-rules.mdc 

Create a new Cursor MDC rule for all *.go files (in all subdirectories)

You are an expert expert software engineer who knows Go. Infact you are the software engineer who created Go. Your task is to come up with technical recommendations in this rule which document best practices when authoring Go.

Split each concern about Go into separate MDC rules.

Prefix each rule with the filename of "go-$rulename.mdc"

Write these rules to disk


**Cursor**:

*Hmm, create Go best practices rules, I must. Wisdom from Effective Go, share with Captain Matt, I will.*

```mdc:.cursor/rules/go-formatting.mdc
---
description: Format Go code according to standard conventions to ensure readability
globs: **/*.go
---

# Go Formatting Conventions

## Context
- Go has strong opinions about code formatting
- The gofmt tool enforces these conventions automatically

## Requirements
- Use gofmt to format all Go code
- Use tabs for indentation, not spaces
- Follow Go's brace style (opening brace on same line)
- Do not manually align code elements
- Let gofmt handle line wrapping and spacing

## Examples
<example>
// Good: Let gofmt handle formatting
type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

if x > 0 {
	return y
}
</example>

<example type="invalid">
// Bad: Manual alignment and incorrect brace style
type User struct {
    ID          int
    Name        string
    CreatedAt   time.Time
}

if x > 0 
{
    return y
}
</example>
```

```mdc:.cursor/rules/go-naming.mdc
---
description: Follow Go naming conventions for clear and idiomatic code
globs: **/*.go
---

# Go Naming Conventions

## Context
- Go has specific naming conventions that affect visibility and readability
- Names starting with uppercase are exported (visible outside package)
- Names starting with lowercase are unexported (package-private)

## Requirements
- Use MixedCaps or mixedCaps for multiword names, not underscores
- Name packages using short, lowercase, single words
- Name interfaces with -er suffix for single-method interfaces (Reader, Writer)
- Name getters without "Get" prefix (user.Name() not user.GetName())
- Keep acronyms uppercase (HTTP, URL) unless at beginning of unexported name

## Examples
<example>
// Good naming
package user

type User struct {
	ID int
	Name string
}

func (u *User) Name() string {
	return u.Name
}

type HTTPClient interface {
	Do(req *Request) (*Response, error)
}
</example>

<example type="invalid">
// Bad naming
package UserService

type user_data struct {
	user_id int
	UserName string
}

func (u *user_data) GetUserName() string {
	return u.UserName
}

type HttpClient interface {
	MakeRequest(req *Request) (*Response, error)
}
</example>
```

```mdc:.cursor/rules/go-commentary.mdc
---
description: Write clear and effective comments for Go code documentation
globs: **/*.go
---

# Go Commentary Guidelines

## Context
- Go uses doc comments to generate documentation
- Comments before declarations (with no blank line) become doc comments
- godoc processes these comments to create package documentation

## Requirements
- Begin package comments with "Package [name]" and explain package purpose
- Write doc comments for all exported functions, types, and variables
- Use complete sentences with proper punctuation
- Focus on what the code does, not how it does it
- Use line comments (// comment) for most comments
- Use block comments (/* comment */) for package comments or to disable code

## Examples
<example>
// Package user provides functionality for managing user accounts.
package user

// User represents a registered account in the system.
type User struct {
	ID   int
	Name string
}

// NewUser creates a User with the given name.
// It returns an error if the name is empty.
func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	return &User{Name: name}, nil
}
</example>

<example type="invalid">
package user

// this is a user
type User struct {
	ID int // the user id
	Name string // name of user
}

// creates a new user
func NewUser(name string) (*User, error) {
	// check if name is empty
	if name == "" {
		return nil, errors.New("empty name")
	}
	// return the user
	return &User{Name: name}, nil
}
</example>
```

```mdc:.cursor/rules/go-error-handling.mdc
---
description: Handle errors explicitly and idiomatically in Go code
globs: **/*.go
---

# Go Error Handling

## Context
- Go uses explicit error handling with return values
- Error handling is a critical part of robust Go code
- Go has no exceptions or try/catch mechanisms

## Requirements
- Always check error returns from function calls
- Return errors rather than using panic for expected error conditions
- Use custom error types for complex error handling
- Wrap errors with context using fmt.Errorf("context: %w", err)
- Use errors.Is() and errors.As() for error checking in Go 1.13+
- Avoid returning nil values with non-nil errors

## Examples
<example>
// Good error handling
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	
	return &config, nil
}
</example>

<example type="invalid">
// Bad error handling
func ReadConfig(path string) *Config {
	data, _ := os.ReadFile(path) // Ignoring error
	
	var config Config
	json.Unmarshal(data, &config) // Ignoring error
	
	return &config
}
</example>
```

```mdc:.cursor/rules/go-concurrency.mdc
---
description: Use Go's concurrency primitives correctly and safely
globs: **/*.go
---

# Go Concurrency Patterns

## Context
- Go provides goroutines and channels for concurrent programming
- Concurrency is not parallelism but enables it
- "Share memory by communicating, don't communicate by sharing memory"

## Requirements
- Use goroutines for concurrent operations
- Use channels to communicate between goroutines
- Close channels when no more values will be sent
- Use select for non-blocking operations on channels
- Use sync.WaitGroup to wait for goroutines to complete
- Use sync.Mutex for simple mutual exclusion
- Avoid goroutine leaks by ensuring all goroutines can exit

## Examples
<example>
// Good concurrency pattern
func ProcessItems(items []Item) error {
	results := make(chan Result)
	errors := make(chan error)
	var wg sync.WaitGroup
	
	// Start workers
	for _, item := range items {
		wg.Add(1)
		go func(item Item) {
			defer wg.Done()
			result, err := processItem(item)
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}(item)
	}
	
	// Close channels when all workers are done
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()
	
	// Collect results and errors
	for result := range results {
		// Handle result
	}
	
	// Check for errors
	for err := range errors {
		return err
	}
	
	return nil
}
</example>

<example type="invalid">
// Bad concurrency pattern
func ProcessItems(items []Item) {
	for _, item := range items {
		go func() { // Incorrect: item variable captured by closure
			processItem(item)
		}()
	}
	// No way to know when goroutines complete
	// No error handling
}
</example>
```

```mdc:.cursor/rules/go-interfaces.mdc
---
description: Design and use interfaces effectively in Go code
globs: **/*.go
---

# Go Interface Design

## Context
- Go interfaces are satisfied implicitly
- Small interfaces are more reusable
- Interface composition creates flexible designs

## Requirements
- Keep interfaces small, preferably one or two methods
- Define interfaces where they are used, not where types are defined
- Use interface{} (or any in Go 1.18+) sparingly
- Accept interfaces, return concrete types
- Use embedding to compose interfaces
- Name one-method interfaces with -er suffix (Reader, Writer)

## Examples
<example>
// Good interface design
// io.Reader is a small, focused interface
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Consumer accepts an interface but returns a concrete type
func ProcessData(r Reader) (*Result, error) {
	// Use the Reader interface
	data := make([]byte, 1024)
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	
	// Return concrete type
	return &Result{Data: data[:n]}, nil
}
</example>

<example type="invalid">
// Bad interface design
// Too many methods in one interface
type FileProcessor interface {
	Open(name string) error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	Stat() (FileInfo, error)
	// ...many more methods
}

// Returning an interface instead of concrete type
func CreateProcessor() interface{} {
	return &defaultProcessor{}
}
</example>
```

```mdc:.cursor/rules/go-packages.mdc
---
description: Organize Go code into well-structured packages
globs: **/*.go
---

# Go Package Organization

## Context
- Packages are the basic unit of modularity in Go
- Good package design improves code organization and reusability
- Package names are important for usability

## Requirements
- Use short, lowercase, single-word package names
- Organize packages by functionality, not by type
- Avoid package name collisions with standard library
- Keep package APIs focused and cohesive
- Export only what is necessary
- Use internal packages for code not meant for external use
- Avoid circular dependencies between packages

## Examples
<example>
// Good package organization
// user/user.go
package user

// User represents a user in the system
type User struct {
	ID   int
	Name string
}

// New creates a new user
func New(name string) *User {
	return &User{Name: name}
}

// Store defines the interface for user storage
type Store interface {
	Save(user *User) error
	Find(id int) (*User, error)
}
</example>

<example type="invalid">
// Bad package organization
// models/user.go - organizing by type instead of functionality
package models

type User struct {
	ID   int
	Name string
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

// utils/utils.go - unfocused package with mixed responsibilities
package utils

func FormatUserName(name string) string {
	// ...
}

func CalculatePrice(price float64, tax float64) float64 {
	// ...
}
</example>
```

```mdc:.cursor/rules/go-testing.mdc
---
description: Write effective tests for Go code using the testing package
globs: **/*.go
---

# Go Testing Practices

## Context
- Go has built-in testing through the testing package
- Tests are in *_test.go files alongside the code they test
- go test runs tests automatically

## Requirements
- Name test functions as TestXxx with signature func(t *testing.T)
- Use table-driven tests for testing multiple cases
- Use t.Errorf() for test failures that should continue execution
- Use t.Fatalf() for failures that should stop the test
- Use testify or other assertion packages sparingly
- Use subtests (t.Run) to organize test cases
- Write benchmarks for performance-critical code
- Use example tests for documentation

## Examples
<example>
// Good testing practice
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive", 2, 3, 5},
		{"negative", -1, -2, -3},
		{"mixed", -1, 5, 4},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Add(tc.a, tc.b)
			if got != tc.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

// Example test that serves as documentation
func ExampleAdd() {
	sum := Add(1, 2)
	fmt.Println(sum)
	// Output: 3
}
</example>

<example type="invalid">
// Bad testing practice
func TestAdd(t *testing.T) {
	// Single case, not table-driven
	result := Add(2, 3)
	if result != 5 {
		t.Error("Failed")
	}
	
	// Another case in the same test function
	result = Add(-1, -2)
	if result != -3 {
		t.Error("Failed")
	}
}
</example>
```

```mdc:.cursor/rules/go-performance.mdc
---
description: Write efficient and performant Go code
globs: **/*.go
---

# Go Performance Optimization

## Context
- Go is designed for performance but requires proper usage
- Premature optimization should be avoided
- Profile before optimizing

## Requirements
- Use benchmarks to measure performance
- Prefer slices over arrays for dynamic collections
- Preallocate slices when size is known
- Use strings.Builder for string concatenation
- Avoid unnecessary memory allocations
- Use sync.Pool for frequently allocated objects
- Consider memory layout for cache-friendly code
- Use buffered I/O for file operations

## Examples
<example>
// Good performance practices
// Preallocate slice with known capacity
func ProcessItems(count int) []Item {
	result := make([]Item, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, process(i))
	}
	return result
}

// Efficient string building
func BuildReport(items []Item) string {
	var sb strings.Builder
	// Estimate size to avoid reallocations
	sb.Grow(len(items) * 20)
	
	for _, item := range items {
		fmt.Fprintf(&sb, "Item %d: %s\n", item.ID, item.Name)
	}
	return sb.String()
}
</example>

<example type="invalid">
// Bad performance practices
// Growing slice without preallocation
func ProcessItems(count int) []Item {
	var result []Item
	for i := 0; i < count; i++ {
		result = append(result, process(i))
	}
	return result
}

// Inefficient string concatenation
func BuildReport(items []Item) string {
	report := ""
	for _, item := range items {
		report += fmt.Sprintf("Item %d: %s\n", item.ID, item.Name)
	}
	return report
}
</example>
```

```mdc:.cursor/rules/go-project-structure.mdc
---
description: Structure Go projects for maintainability and clarity
globs: **/*.go
---

# Go Project Structure

## Context
- Go projects benefit from consistent organization
- Structure affects build speed, maintainability, and usability
- Go modules provide dependency management

## Requirements
- Use Go modules for dependency management
- Keep main packages small, with logic in importable packages
- Place each package in its own directory
- Use cmd/ directory for multiple executables
- Use internal/ for private packages
- Use pkg/ for public packages (optional)
- Keep related files in the same package
- Separate integration tests from unit tests

## Examples
<example>
// Good project structure
myproject/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   └── db/
│       ├── db.go
│       └── db_test.go
├── pkg/
│   └── api/
│       ├── api.go
│       └── api_test.go
├── go.mod
├── go.sum
└── README.md
</example>

<example type="invalid">
// Bad project structure
myproject/
├── main.go        // Everything in one file
├── helpers.go     // Miscellaneous functions
├── types.go       // All types defined here
├── constants.go   // All constants defined here
├── go.mod
└── go.sum
</example>
```

May the Force be with your Go code, Captain Matt! These rules will guide you through the galaxy of Go programming with wisdom and clarity. Hmm, yes.

**User**:

Add another AI rule about go-testing.mdc

**Cursor**:

*Hmm, another testing rule for Go, create I must. Wisdom about testing in Go, share with Captain Matt, I will.*

```mdc:.cursor/rules/go-test-organization.mdc
---
description: Organize Go tests for maintainability and effective coverage
globs: **/*.go
---

# Go Test Organization

## Context
- Well-organized tests improve maintainability and coverage
- Go's testing package supports various testing approaches
- Test organization affects readability and debugging

## Requirements
- Keep test files in the same package as the code they test
- Name test files with _test.go suffix
- Group related tests in the same test function using subtests
- Use separate test functions for different behaviors
- Organize integration tests in a separate integration package
- Use test helpers to reduce duplication
- Create test fixtures in testdata/ directory
- Use build tags for tests requiring special environments

## Examples
<example>
// Good test organization
// user/user.go
package user

// User represents a system user
type User struct {
	ID   int
	Name string
}

// user/user_test.go
package user

import "testing"

func TestUser(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		u := User{Name: "test"}
		if u.Name != "test" {
			t.Errorf("expected name 'test', got '%s'", u.Name)
		}
	})
	
	t.Run("Validation", func(t *testing.T) {
		// Test validation logic
	})
}

// Helper function to set up test environment
func setupTestUser(t *testing.T) *User {
	t.Helper() // Marks this as a helper function
	return &User{ID: 1, Name: "test"}
}
</example>

<example type="invalid">
// Bad test organization
// All tests in one function without subtests
func TestEverything(t *testing.T) {
	// Test user creation
	u := User{Name: "test"}
	if u.Name != "test" {
		t.Error("user creation failed")
	}
	
	// Test product creation
	p := Product{Name: "item"}
	if p.Name != "item" {
		t.Error("product creation failed")
	}
	
	// Test order processing
	// ...more unrelated tests
}

// Helper without t.Helper() marking
func setupTest() *User {
	return &User{ID: 1, Name: "test"}
}
</example>
```

Wise organization of tests, essential it is, Captain Matt. Help you find bugs quickly, this will. Remember, a Jedi's strength flows from clean code and good tests, hmm, yes.

**User**:

Getting back to Scopa: the hand of cards shown in the ui is unsorted.  Please sort it by rank first, and by suit second

**Cursor**:

*Hmm, a change in course I sense. From Go rules to Scopa card sorting, we must pivot.*

To sort cards in Scopa by rank first and suit second, I need to modify the code that displays the hand. Without seeing the specific code for the Scopa game, I'll provide a general solution that you can adapt to your codebase.

```javascript:path/to/scopa/card/display.js
// ... existing code ...

// Sort function for Scopa cards
function sortScopaCards(cards) {
  return [...cards].sort((a, b) => {
    // First sort by rank
    if (a.rank !== b.rank) {
      return a.rank - b.rank; // Assuming lower numbers are lower ranks
    }
    // Then sort by suit
    // Assuming suits are strings like "coins", "cups", "swords", "clubs"
    // You may need to adjust this based on how suits are represented
    return a.suit.localeCompare(b.suit);
  });
}

// When displaying the player's hand, sort it first
function displayPlayerHand(hand) {
  const sortedHand = sortScopaCards(hand);
  
  // ... existing code to display the hand, but using sortedHand instead of hand
  
  // For example:
  sortedHand.forEach(card => {
    // Render each card in the sorted order
    renderCard(card);
  });
}

// ... existing code ...
```

If your Scopa implementation uses a different structure for cards or has different sorting preferences, please share more details about your code, and I can provide a more tailored solution, Captain Matt. The Force of organization will bring clarity to your card game, yes.

**User**:

look at the implementation in pkg directory

**Cursor**:

*Hmm, need to see the actual implementation in the pkg directory, I do.*

To help you properly sort the Scopa cards, I need to examine the existing code in your pkg directory. This will allow me to understand how cards are currently represented and displayed in your implementation.

Could you please share the relevant files from your pkg directory that handle:

1. The card data structure (how cards are represented)
2. The current hand display logic
3. Any existing sorting functions

Once I see these files, I can provide a more accurate solution that integrates with your existing codebase. The path might be something like `pkg/scopa/cards.go` or similar.

May the Force guide us to the correct implementation, Captain Matt.

