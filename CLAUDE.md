# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go code generation tool (`github.com/kmio11/codegen`) that provides two main capabilities:
1. **Mock generation** from Go interfaces - creates mock implementations and stub helpers for testing
2. **Interface generation** from Go structs - extracts interface definitions from struct methods for clean architecture

## Development Commands

### Building and Running
```bash
# Build the tool
go build .

# Install globally
go get github.com/kmio11/codegen

# Run mock generation
go run . mock -pkg <package> -type <interface> -out <output_file>

# Run interface generation
go run . interface -pkg <package> -type <struct> -out <output_file>

# Example usage (from examples directory)
cd _examples/mock
go run ../.. mock -pkg . -type SomeInterface -out testing_gen.go

cd _examples/interface  
go run ../.. interface -pkg . -type UserService -out user_service_gen.go
```

### Quality Checks
```bash
# Run gofmt to format code
gofmt -w .
# Run go vet to check for issues
go vet ./...
# Run golangci-lint for linting
staticcheck ./...
```

### Testing
```bash
# Run tests (standard Go testing)
go test ./...

# Run sample tests
cd cmd/mock/_sample && go test
```

## Code Architecture

### Command Pattern CLI
- Main entry: `main.go` with pluggable command system
- Commands implement `Command` interface: `Name()`, `Description()`, `Usage()`, `Parse()`, `Execute()`
- Currently supports two commands: `mock` and `interface`

### Three-Layer Generation Architecture

1. **Parser Layer** (`generator/parser/`)
   - Uses `golang.org/x/tools` to parse Go packages and extract interface/struct definitions
   - Converts Go AST into internal model representations
   - Supports both interface parsing (for mocking) and struct parsing (for interface generation)

2. **Model Layer** (`generator/model/`)
   - Intermediate representation: Package, File, Interface, Method, Type structures
   - Handles import management and dependency resolution

3. **Generator Layer** (`generator/generator.go`)
   - Generates formatted Go source code from models
   - Uses `go/format` for proper code formatting
   - Provides fluent API for code generation

### Mock Generation Strategy
For each interface, generates:
- **Mock struct**: Contains function fields (`FakeMethodName`) that can be set to control behavior
- **Stub struct**: Provides convenient way to set return values with `NewMock()` method

## Key Files and Components

- `main.go` - Main CLI dispatcher and command registration
- `cmd/mock/mock.go` - Mock generation command implementation
- `cmd/interface/interface.go` - Interface generation command implementation
- `generator/generator.go` - Core code generation with fluent API
- `generator/parser/` - Go package parsing and AST processing (supports both interfaces and structs)
- `generator/model/` - Internal code representation models
- `_examples/mock/` - Mock generation example usage and test files
- `_examples/interface/` - Interface generation example usage and test files

## Dependencies

- `golang.org/x/tools` - Required for Go code analysis and manipulation
- Standard library only for core functionality

## Usage Patterns

### Interface Generation Pattern

Extract interfaces from struct implementations:
```go
// Input: Struct with methods
type UserService struct { db Database }
func (u *UserService) GetUser(id string) (*User, error) { ... }
func (u *UserService) CreateUser(user *User) error { ... }

// Generated: Clean interface definition
type UserServiceInterface interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}
```

### Mock Usage Pattern

Generated mocks embed the original interface and provide:
```go
// Generated structures
type MockSomeInterface struct {
    SomeInterface           // Embedded interface
    FakeMethodName func(...) // Function fields for each method
}

type StubSomeInterface struct {
    MethodName StubMethodName // Stub for each method
}

// Usage in tests
stub := StubSomeInterface{MethodName: StubMethodName{R0: "result"}}
mock := stub.NewMock() // Returns interface implementation
```

## TDD Development Workflow

Follow this Test-Driven Development workflow for all code changes:

### 1. Red Phase - Write Failing Tests
```bash
# Create tests for new functionality
# Test files: *_test.go
go test ./... # Confirm tests fail
```

### 2. Green Phase - Minimal Implementation
```bash
# Write minimal code to make tests pass
go test ./... # Confirm tests pass
```

### 3. Refactor Phase - Improve Code
```bash
# Refactor code while keeping tests green
go test ./... # Confirm tests still pass
```

### 4. Cleanup Phase - Remove Unnecessary Code
```bash
# Remove unnecessary comments, debug code, unused variables
# Remove mock code, temporary test code
# Remove unused imports
go mod tidy # Remove unused dependencies
```

### 5. Quality Checks - Code Quality
```bash
# Format code
gofmt -w .

# Static analysis
go vet ./...
staticcheck ./...

# Run all tests
go test ./...
```

### 6. Final Verification
```bash
# Final check - All Green confirmation
go test ./... && echo "✅ All tests passed"
```

### Development Cycle
1. **Red** → Write failing tests
2. **Green** → Minimal implementation  
3. **Refactor** → Improve code structure
4. **Cleanup** → Remove unnecessary code
5. **Quality** → Run quality checks
6. **Repeat** → Next feature

### Cleanup Targets (Unnecessary Code Examples)
- Debug print/log statements
- Commented out old code
- Unused variables, functions, imports
- Temporary mock code for testing
- Completed TODO comments

This workflow ensures clean, maintainable, and well-tested code.

## Go Coding Rules

Follow these Go coding conventions and best practices:

### Naming and Style
- Write idiomatic, readable Go code following community conventions
- Use MixedCaps naming; avoid package name duplication and excessive abbreviations
- Document all exported symbols with clear Go documentation comments

### Error Handling
- Handle errors explicitly with error type; use panic only for unexpected runtime errors
- Use early returns for error handling to keep normal code path unindented

### Architecture Patterns
- Avoid global state in libraries; use instances and explicit dependency injection
- Pass context.Context as first function parameter; never store in struct fields

### Modern Go Features
- Use `any` instead of `interface{}` (since Go 1.18)