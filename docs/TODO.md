# Struct-to-Interface Generation Command Implementation TODO

## Overview
✅ **COMPLETED** - Implement a new command `interface` that generates Go interfaces from struct types, extracting all exported methods and creating a corresponding interface definition.

## Command Specification
- **Command name**: `interface`
- **Purpose**: Generate interface from struct methods
- **Usage**: `go run ./cmd interface -pkg <package> -type <struct> -out <output_file>`

## Implementation Checklist

### Phase 1: Parser Layer Extension - ✅ COMPLETED
- [x] **Extend parser to support struct parsing**
  - [x] Add `isStruct()` function in `generator/parser/parser.go`
  - [x] Implement `parseStructObj()` method to extract struct methods
  - [x] Update `setContents()` method to handle struct types (line 87-95)
  - [x] Use `types.MethodSet` to get all methods including pointer receiver methods
  - [x] Filter exported methods only (capitalize check)

- [x] **Add struct method extraction logic**  
  - [x] Create method extraction in `generator/parser/object.go`
  - [x] Handle method signatures (parameters, return types)
  - [x] Support generic type parameters in struct methods
  - [x] Preserve method documentation comments if available

### Phase 2: Command Implementation - ✅ COMPLETED
- [x] **Create command structure**
  - [x] Create `cmd/interface/` directory
  - [x] Implement `cmd/interface/interface.go` with Command interface
  - [x] Add flag definitions: `-pkg`, `-type`, `-out`, `-outpkg`, `-selfpkg`, `-name`
  - [x] Follow existing mock command pattern from `cmd/mock/mock.go`

- [x] **Implement command logic**
  - [x] Parse command line flags
  - [x] Validate required flags (type is mandatory)
  - [x] Set default values (name defaults to `<StructName>Interface`)
  - [x] Handle file output vs stdout
  - [x] Add proper error handling and user feedback

- [x] **Register command**
  - [x] Add import in `cmd/codegen.go`
  - [x] Add to `commands` slice for registration
  - [x] Update help text to include new command

### Phase 3: Interface Generation Logic - ✅ COMPLETED
- [x] **Implement interface generation**
  - [x] Create interface model from extracted struct methods
  - [x] Generate interface name (default: `<StructName>Interface`)
  - [x] Handle generic type parameters correctly
  - [x] Generate proper method signatures without receiver

- [x] **Code generation integration**
  - [x] Use existing `generator.Generator` fluent API
  - [x] Generate proper package declaration
  - [x] Add import statements for referenced types
  - [x] Format generated code using `go/format`
  - [x] Add generation attribution header

### Phase 4: Testing and Validation - ✅ COMPLETED
- [x] **Unit Tests**
  - [x] Test struct method extraction
  - [x] Test interface generation logic
  - [x] Test flag parsing and validation
  - [x] Test error handling scenarios

- [x] **Integration Tests**
  - [x] Create `cmd/interface/_sample/` directory
  - [x] Add sample struct with various method signatures
  - [x] Test basic struct to interface generation
  - [x] Test generic struct to interface generation
  - [x] Test complex method signatures (multiple params, returns)

- [x] **Edge Case Testing**
  - [x] Struct with no exported methods
  - [x] Struct with only private methods
  - [x] Struct with embedded interfaces
  - [x] Struct with generic type parameters
  - [x] Non-existent struct name
  - [x] Invalid package path

### Phase 5: Documentation and Examples - ✅ COMPLETED
- [x] **Usage Examples**
  - [x] Add usage examples to command help text
  - [x] Create example in `cmd/interface/_sample/`
  - [x] Document common use cases

- [x] **Command Help**
  - [x] Implement detailed `Usage()` method
  - [x] Add flag descriptions
  - [x] Add examples to help output

### Phase 6: Quality Assurance - ✅ COMPLETED
- [x] **Code Quality**
  - [x] Run code formatting: `gofmt -w .`
  - [x] Run static analysis: `go vet ./...`, `staticcheck ./...`
  - [x] Follow existing code conventions and patterns
  - [x] Add proper error messages and logging

- [x] **Testing**
  - [x] Run all tests: `go test ./...`
  - [x] Test command manually with various inputs
  - [x] Verify generated interface code compiles
  - [x] Test with real-world struct examples

### Phase 7: Final Integration - ✅ COMPLETED
- [x] **End-to-End Testing**
  - [x] Test complete workflow: struct → interface → compilation
  - [x] Verify interface can be implemented by original struct
  - [x] Test with various project structures

- [x] **Performance Testing**
  - [x] Test with large structs (many methods)
  - [x] Test with complex generic constraints
  - [x] Verify reasonable performance

## Technical Implementation Notes

### Parser Extension Details
```go
// In generator/parser/parser.go, around line 87-95:
if types.IsInterface(obj.Type()) {
    // existing interface logic
} else if isStruct(obj.Type()) {  // NEW
    intf, err := p.parseStructAsInterface(obj)
    if err != nil {
        return err
    }
    pkg.Interfaces = append(pkg.Interfaces, intf)
} else {
    return fmt.Errorf("%s is unsupported", obj.Type())
}
```

### Method Extraction Logic
```go
func (p *Parser) parseStructAsInterface(obj types.Object) (*model.Interface, error) {
    // Get method set including pointer receiver methods
    methodSet := types.NewMethodSet(types.NewPointer(obj.Type()))
    
    // Extract exported methods only
    var methods []*model.Method
    for i := 0; i < methodSet.Len(); i++ {
        method := methodSet.At(i)
        if method.Obj().Exported() {
            // Convert to model.Method
        }
    }
    
    return model.NewInterface(interfaceName, methods), nil
}
```

### Command Flag Structure
```go
type Command struct {
    fs              *flag.FlagSet
    flagPkg         *string  // package path
    flagType        *string  // struct name (required)
    flagOut         *string  // output file
    flagOutPkg      *string  // output package name
    flagSelfPkgPath *string  // self package path
    flagName        *string  // interface name (optional)
}
```

## Success Criteria - ✅ ALL COMPLETED
- [x] Command successfully generates interfaces from structs
- [x] Generated interfaces compile without errors  
- [x] Original struct can implement generated interface
- [x] All tests pass
- [x] Code follows project conventions
- [x] Documentation is complete and accurate

## Example Usage
```bash
# Basic usage
go run ./cmd interface -pkg . -type UserService -out user_interface.go

# With custom interface name
go run ./cmd interface -pkg ./internal/service -type UserService -name UserServiceInterface -out user_interface.go

# Output to different package
go run ./cmd interface -pkg ./internal -type Handler -outpkg contracts -out ./contracts/handler.go
```

## Generated Output Example
```go
// Generated from struct UserService
type UserServiceInterface interface {
    CreateUser(name string, email string) (*User, error)
    GetUser(id int64) (*User, error)
    UpdateUser(id int64, user *User) error
    DeleteUser(id int64) error
}
```