package parser

import (
	"go/types"
	"log"
	"os"
	"testing"
)

func TestNewParser(t *testing.T) {
	// Test parser creation with default options
	parser := NewParser()

	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}

	if parser.log == nil {
		t.Error("NewParser() should have default logger")
	}

	// Check default values
	if parser.ParsedPkg != nil {
		t.Error("ParsedPkg should be nil by default")
	}

	if parser.Targets != nil {
		t.Error("Targets should be nil by default")
	}

	if parser.stopLoadErr {
		t.Error("stopLoadErr should be false by default")
	}
}

func TestNewParserWithOptions(t *testing.T) {
	logger := log.New(os.Stderr, "TEST: ", log.LstdFlags)
	targets := []string{"TestInterface"}

	parser := NewParser(
		OptLogger(logger),
		OptParseTarget(targets),
		OptStopLoadErr(),
	)

	if parser == nil {
		t.Fatal("NewParser() with options returned nil")
	}

	if parser.log != logger {
		t.Error("NewParser() should use provided logger")
	}

	if len(parser.Targets) != 1 || parser.Targets[0] != "TestInterface" {
		t.Errorf("NewParser() targets = %v, want %v", parser.Targets, targets)
	}

	if !parser.stopLoadErr {
		t.Error("stopLoadErr should be true when OptStopLoadErr is used")
	}
}

func TestParserOptionsApplication(t *testing.T) {
	tests := []struct {
		name     string
		option   Opts
		validate func(*Parser) bool
	}{
		{
			name:   "logger option",
			option: OptLogger(log.New(os.Stdout, "TEST: ", log.LstdFlags)),
			validate: func(p *Parser) bool {
				return p.log != nil
			},
		},
		{
			name:   "parse target option",
			option: OptParseTarget([]string{"Interface1", "Interface2"}),
			validate: func(p *Parser) bool {
				return len(p.Targets) == 2 && p.Targets[0] == "Interface1" && p.Targets[1] == "Interface2"
			},
		},
		{
			name:   "stop load error option",
			option: OptStopLoadErr(),
			validate: func(p *Parser) bool {
				return p.stopLoadErr == true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.option)

			if !tt.validate(parser) {
				t.Errorf("Option %s not applied correctly", tt.name)
			}
		})
	}
}

func TestParserNilLoggerHandling(t *testing.T) {
	parser := NewParser(OptLogger(nil))

	if parser.log == nil {
		t.Error("Parser should have a logger even when nil is provided")
	}

	// Test that we can write to the logger without panic
	parser.log.Println("test message")
}

func TestParserParseMethod(t *testing.T) {
	parser := NewParser()

	// Test with nil ParsedPkg
	_, err := parser.Parse()
	if err == nil {
		t.Error("Parse() should return error when ParsedPkg is nil")
	}

	// Test with invalid target count
	parser.ParsedPkg = &Package{
		Pkg: types.NewPackage("test", "test"),
	}
	parser.Targets = []string{"Interface1", "Interface2"} // More than 1

	_, err = parser.Parse()
	if err == nil {
		t.Error("Parse() should return error when target count is not 1")
	}

	// Test with no targets
	parser.Targets = []string{}
	_, err = parser.Parse()
	if err == nil {
		t.Error("Parse() should return error when target count is 0")
	}
}

func TestNewTypeParser(t *testing.T) {
	parser := NewParser()
	tp := parser.newTypeParser()

	if tp == nil {
		t.Fatal("newTypeParser() returned nil")
	}

	if tp.p != parser {
		t.Error("newTypeParser() should reference parent parser")
	}

	if tp.parsed == nil {
		t.Error("newTypeParser() should initialize parsed map")
	}

	if tp.stats == nil {
		t.Error("newTypeParser() should initialize stats map")
	}
}

func TestPackageCreation(t *testing.T) {
	pkg := types.NewPackage("example.com/test", "test")
	parserPkg := &Package{
		Pkg: pkg,
	}

	if parserPkg.Pkg.Path() != "example.com/test" {
		t.Errorf("Package path = %v, want %v", parserPkg.Pkg.Path(), "example.com/test")
	}

	if parserPkg.Pkg.Name() != "test" {
		t.Errorf("Package name = %v, want %v", parserPkg.Pkg.Name(), "test")
	}
}

func TestParseMultipleOptions(t *testing.T) {
	logger := log.New(os.Stderr, "MULTI: ", log.LstdFlags)
	targets := []string{"TestInterface"}

	parser := NewParser(
		OptLogger(logger),
		OptParseTarget(targets),
		OptStopLoadErr(),
	)

	// All options should be applied
	if parser.log != logger {
		t.Error("Logger option not applied")
	}

	if len(parser.Targets) != 1 || parser.Targets[0] != "TestInterface" {
		t.Error("Target option not applied")
	}

	if !parser.stopLoadErr {
		t.Error("StopLoadErr option not applied")
	}
}

func TestOptPackageOption(t *testing.T) {
	pkg := &Package{
		Pkg: types.NewPackage("test", "test"),
	}

	parser := NewParser(OptPackage(pkg))

	if parser.ParsedPkg != pkg {
		t.Error("OptPackage should set ParsedPkg")
	}
}

func TestParserState(t *testing.T) {
	parser := NewParser()

	// Test initial state
	if parser.ParsedPkg != nil {
		t.Error("Initial ParsedPkg should be nil")
	}

	if parser.Targets != nil {
		t.Error("Initial Targets should be nil")
	}

	if parser.stopLoadErr {
		t.Error("Initial stopLoadErr should be false")
	}
}

// Test struct parsing support
func TestStructParsing(t *testing.T) {
	parser := NewParser()

	// Create a mock package with a struct
	pkg := types.NewPackage("test", "test")
	scope := pkg.Scope()

	// Create a struct type
	structType := types.NewNamed(
		types.NewTypeName(0, pkg, "TestStruct", nil),
		types.NewStruct([]*types.Var{
			types.NewVar(0, pkg, "field", types.Typ[types.String]),
		}, nil),
		nil,
	)

	// Create methods for the struct
	// Method 1: Name() string
	nameMethod := types.NewFunc(0, pkg, "Name", types.NewSignatureType(
		nil, // receiver
		nil, // recv type params
		nil, // type params
		nil, // params
		types.NewTuple(types.NewVar(0, pkg, "", types.Typ[types.String])), // results
		false, // variadic
	))

	// Method 2: SetName(name string)
	setNameMethod := types.NewFunc(0, pkg, "SetName", types.NewSignatureType(
		nil, // receiver
		nil, // recv type params
		nil, // type params
		types.NewTuple(types.NewVar(0, pkg, "name", types.Typ[types.String])), // params
		nil,   // results
		false, // variadic
	))

	// Add methods to the struct type
	structType.AddMethod(nameMethod)
	structType.AddMethod(setNameMethod)

	// Add struct to package scope
	scope.Insert(structType.Obj())

	parser.ParsedPkg = &Package{Pkg: pkg}
	parser.Targets = []string{"TestStruct"}

	// Parse the struct
	modelPkg, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}

	// Verify interface was created
	if len(modelPkg.Interfaces) != 1 {
		t.Errorf("Expected 1 interface, got %d", len(modelPkg.Interfaces))
	}

	intf := modelPkg.Interfaces[0]
	if intf.Name() != "TestStructInterface" {
		t.Errorf("Expected interface name 'TestStructInterface', got %s", intf.Name())
	}

	// Verify methods were extracted
	methods := intf.Methods()
	if len(methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(methods))
	}

	// Check method names
	methodNames := make(map[string]bool)
	for _, method := range methods {
		methodNames[method.Name()] = true
	}

	if !methodNames["Name"] {
		t.Error("Expected method 'Name' not found")
	}
	if !methodNames["SetName"] {
		t.Error("Expected method 'SetName' not found")
	}
}
