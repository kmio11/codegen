# codegen

A Go code generation tool that creates mock objects from Go interfaces with full support for generics.

## Installation

```bash
# Install globally
go get github.com/kmio11/codegen/cmd

# Or clone and build locally
git clone https://github.com/kmio11/codegen
cd codegen
go build ./cmd
```

## Quick Start

### Basic Interface Mocking

Given an interface:
```go
// user.go
type UserService interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
    DeleteUser(id string) error
}
```

Generate mocks:
```bash
go run github.com/kmio11/codegen/cmd mock -pkg . -type UserService -out user_mock.go
```

This generates two types of mock structures:

**Mock struct (with function fields for direct control):**
```go
type MockUserService struct {
    UserService
    FakeGetUser    func(id string) (*User, error)
    FakeCreateUser func(user *User) error
    FakeDeleteUser func(id string) error
}
```

**Stub struct (convenient for testing):**
```go
type StubUserService struct {
    GetUser    StubGetUser
    CreateUser StubCreateUser
    DeleteUser StubDeleteUser
}

func (s StubUserService) NewMock() UserService {
    return &MockUserService{
        FakeGetUser:    func(a0 string) (*User, error) { return s.GetUser.R0, s.GetUser.R1 },
        FakeCreateUser: func(a0 *User) error { return s.CreateUser.R0 },
        FakeDeleteUser: func(a0 string) error { return s.DeleteUser.R0 },
    }
}
```

### Generic Interface Mocking

For generic interfaces:
```go
// repository.go
type Repository[T any] interface {
    Get(id string) (T, error)
    Save(item T) error
    List() ([]T, error)
}

type Cache[K comparable, V any] interface {
    Set(key K, value V)
    Get(key K) (V, bool)
    Delete(key K)
}
```

Generate with type parameters:
```bash
go run github.com/kmio11/codegen/cmd mock -pkg . -pkg . -type Repository -out repository_mock_gen.go
go run github.com/kmio11/codegen/cmd mock -pkg . -type "Cache" -out cache_mock_gen.go
```

## Usage Examples

### Direct Mock Usage
```go
func TestUserService_Direct(t *testing.T) {
    mock := &MockUserService{
        FakeGetUser: func(id string) (*User, error) {
            if id == "123" {
                return &User{ID: "123", Name: "John"}, nil
            }
            return nil, errors.New("user not found")
        },
    }
    
    user, err := mock.GetUser("123")
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### Stub-Based Testing (Recommended)
```go
func TestUserService_Stub(t *testing.T) {
    stub := StubUserService{
        GetUser: StubGetUser{
            R0: &User{ID: "123", Name: "John"},
            R1: nil,
        },
        CreateUser: StubCreateUser{R0: nil},
    }
    
    userService := stub.NewMock()
    
    user, err := userService.GetUser("123")
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### Generic Mock Usage
```go
func TestRepository_Generic(t *testing.T) {
    // Create mock for Repository[string]
    stub := StubRepository[string]{
        Get: StubGet[string]{
            R0: "test-item",
            R1: nil,
        },
        Save: StubSave[string]{R0: nil},
    }
    
    repo := stub.NewMock()
    
    item, err := repo.Get("key")
    assert.NoError(t, err)
    assert.Equal(t, "test-item", item)
}
```

## Command Reference

### Mock Command

```bash
go run github.com/kmio11/codegen/cmd mock [options]
```

**Required Options:**
- `-pkg <package>` - Target package path
- `-type <interface>` - Interface name to mock
- `-out <file>` - Output file path

**Optional Options:**
- `-outpkg <package>` - Output package name
- `-selfpkgpath <path>` - Self package path for imports

**Examples:**
```bash
# Basic usage
go run ./cmd mock -pkg . -type UserService -out user_mock.go

# With custom output package
go run ./cmd mock -pkg ./services -type PaymentService -outpkg mocks -out mocks/payment_mock.go

# Generic interface
go run ./cmd mock -pkg . -type "Repository[T any]" -out repository_mock.go

# Interface with complex constraints  
go run ./cmd mock -pkg . -type "Processor[T comparable, U int|float64]" -out processor_mock.go
```

## Features

- ✅ **Go 1.18+ Generics Support** - Generate mocks for generic interfaces with type parameters and constraints
- ✅ **Dual Mock Strategy** - Creates both Mock structs (function fields) and Stub structs (convenient testing)
- ✅ **Smart Import Management** - Automatic dependency resolution and package import handling
- ✅ **Clean Code Generation** - Produces properly formatted, idiomatic Go code

## Examples

The `cmd/mock/_sample/` directory contains working examples:

- **`sample.go`** - Basic interface definitions
- **`generic_interfaces.go`** - Generic interface examples with various constraints
- **`*_test.go`** files - Demonstrate proper mock usage patterns
- **Generated mock files** - Show expected output format