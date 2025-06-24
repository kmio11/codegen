package parser

import (
	"go/types"
	"log"
	"os"
	"testing"
)

func TestParseInterfaceObj(t *testing.T) {
	// Create a simple parser for testing
	parser := NewParser(
		OptLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)),
	)

	tests := []struct {
		name           string
		setupInterface func() types.Object
		expectGeneric  bool
		expectMethods  int
	}{
		{
			name: "non-generic interface",
			setupInterface: func() types.Object {
				// Create a simple interface type for testing
				methods := []*types.Func{}
				
				// Create method: Process(data string) error
				sig := types.NewSignatureType(
					nil, nil, nil,
					types.NewTuple(
						types.NewParam(0, nil, "data", types.Typ[types.String]),
					),
					types.NewTuple(
						types.NewParam(0, nil, "", types.Universe.Lookup("error").Type()),
					),
					false,
				)
				method := types.NewFunc(0, nil, "Process", sig)
				methods = append(methods, method)
				
				// Create interface type
				interfaceType := types.NewInterfaceType(methods, nil)
				
				// Create a package and type name
				pkg := types.NewPackage("example.com/test", "test")
				typeName := types.NewTypeName(0, pkg, "Processor", interfaceType)
				
				return typeName
			},
			expectGeneric: false,
			expectMethods: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := tt.setupInterface()
			if obj == nil {
				t.Skip("Skipping test - interface setup failed")
				return
			}
			
			result, err := parser.parseInterfaceObj(obj)
			if err != nil {
				t.Fatalf("parseInterfaceObj() error = %v", err)
			}
			
			if result == nil {
				t.Fatal("parseInterfaceObj() returned nil")
			}
			
			if result.IsGeneric() != tt.expectGeneric {
				t.Errorf("parseInterfaceObj() IsGeneric() = %v, want %v", result.IsGeneric(), tt.expectGeneric)
			}
			
			if len(result.Methods()) != tt.expectMethods {
				t.Errorf("parseInterfaceObj() Methods() count = %v, want %v", len(result.Methods()), tt.expectMethods)
			}
		})
	}
}

func TestGetMethodSet(t *testing.T) {
	parser := &Parser{}
	
	// Create a simple interface for testing
	methods := []*types.Func{}
	
	// Create method: Get() string
	sig := types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(), // no parameters
		types.NewTuple(
			types.NewParam(0, nil, "", types.Typ[types.String]),
		),
		false,
	)
	method := types.NewFunc(0, nil, "Get", sig)
	methods = append(methods, method)
	
	// Create interface type
	interfaceType := types.NewInterfaceType(methods, nil)
	
	// Create a package and type name
	pkg := types.NewPackage("example.com/test", "test")
	typeName := types.NewTypeName(0, pkg, "Getter", interfaceType)
	
	// Test getMethodSet
	methodSet, err := parser.getMethodSet(typeName, false)
	if err != nil {
		t.Fatalf("getMethodSet() error = %v", err)
	}
	
	if methodSet == nil {
		t.Fatal("getMethodSet() returned nil")
	}
	
	if methodSet.Len() != 1 {
		t.Errorf("getMethodSet() length = %v, want %v", methodSet.Len(), 1)
	}
	
	// Test with pointer
	methodSetPtr, err := parser.getMethodSet(typeName, true)
	if err != nil {
		t.Fatalf("getMethodSet() with pointer error = %v", err)
	}
	
	if methodSetPtr == nil {
		t.Fatal("getMethodSet() with pointer returned nil")
	}
}

func TestParseInterfaceObjWithComplexMethods(t *testing.T) {
	parser := NewParser(
		OptLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)),
	)

	// Create interface with multiple methods
	methods := []*types.Func{}
	
	// Method 1: Save(item string) error
	saveSig := types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(
			types.NewParam(0, nil, "item", types.Typ[types.String]),
		),
		types.NewTuple(
			types.NewParam(0, nil, "", types.Universe.Lookup("error").Type()),
		),
		false,
	)
	saveMethod := types.NewFunc(0, nil, "Save", saveSig)
	methods = append(methods, saveMethod)
	
	// Method 2: Get(id int) (string, error)
	getSig := types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(
			types.NewParam(0, nil, "id", types.Typ[types.Int]),
		),
		types.NewTuple(
			types.NewParam(0, nil, "", types.Typ[types.String]),
			types.NewParam(0, nil, "", types.Universe.Lookup("error").Type()),
		),
		false,
	)
	getMethod := types.NewFunc(0, nil, "Get", getSig)
	methods = append(methods, getMethod)
	
	// Create interface type
	interfaceType := types.NewInterfaceType(methods, nil)
	
	// Create a package and type name
	pkg := types.NewPackage("example.com/test", "test")
	typeName := types.NewTypeName(0, pkg, "Repository", interfaceType)
	
	result, err := parser.parseInterfaceObj(typeName)
	if err != nil {
		t.Fatalf("parseInterfaceObj() error = %v", err)
	}
	
	if result == nil {
		t.Fatal("parseInterfaceObj() returned nil")
	}
	
	if result.Name() != "Repository" {
		t.Errorf("parseInterfaceObj() Name() = %v, want %v", result.Name(), "Repository")
	}
	
	if len(result.Methods()) != 2 {
		t.Errorf("parseInterfaceObj() Methods() count = %v, want %v", len(result.Methods()), 2)
	}
	
	// Check method names
	methodNames := make(map[string]bool)
	for _, method := range result.Methods() {
		methodNames[method.Name()] = true
	}
	
	if !methodNames["Save"] {
		t.Error("parseInterfaceObj() missing Save method")
	}
	
	if !methodNames["Get"] {
		t.Error("parseInterfaceObj() missing Get method")
	}
}