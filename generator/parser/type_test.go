package parser

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/kmio11/codegen/generator/model"
)

func TestParseTypeBasicTypes(t *testing.T) {
	// Create a typeParser instance
	p := NewParser()
	tp := p.newTypeParser()

	// Test basic types parsing
	tests := []struct {
		name     string
		goType   types.Type
		expected string
	}{
		{
			name:     "string type",
			goType:   types.Typ[types.String],
			expected: "string",
		},
		{
			name:     "int type",
			goType:   types.Typ[types.Int],
			expected: "int",
		},
		{
			name:     "bool type",
			goType:   types.Typ[types.Bool],
			expected: "bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.parseType(tt.goType)
			if err != nil {
				t.Fatalf("parseType() error = %v", err)
			}

			if result == nil {
				t.Fatal("parseType() returned nil")
			}

			typeStr := result.PrintType("", model.PackageMap{})
			if typeStr != tt.expected {
				t.Errorf("PrintType() = %v, want %v", typeStr, tt.expected)
			}
		})
	}
}

func TestParseSliceType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create slice type []string
	sliceType := types.NewSlice(types.Typ[types.String])

	result, err := tp.parseType(sliceType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "[]string"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseMapType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create map type map[string]int
	mapType := types.NewMap(types.Typ[types.String], types.Typ[types.Int])

	result, err := tp.parseType(mapType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "map[string]int"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseChannelType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create channel type chan string
	chanType := types.NewChan(types.SendRecv, types.Typ[types.String])

	result, err := tp.parseType(chanType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "chan string"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParsePointerType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create pointer type *string
	ptrType := types.NewPointer(types.Typ[types.String])

	result, err := tp.parseType(ptrType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "*string"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseFunctionType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create function type func(string) int
	params := types.NewTuple(
		types.NewVar(token.NoPos, nil, "s", types.Typ[types.String]),
	)
	results := types.NewTuple(
		types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
	)
	funcType := types.NewSignatureType(nil, nil, nil, params, results, false)

	result, err := tp.parseType(funcType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "func(s string) int"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseInterfaceType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create empty interface{} type
	interfaceType := types.NewInterfaceType(nil, nil)

	result, err := tp.parseType(interfaceType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "interface{}"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseStructType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create struct type with fields
	fields := []*types.Var{
		types.NewField(token.NoPos, nil, "Name", types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "Age", types.Typ[types.Int], false),
	}
	structType := types.NewStruct(fields, nil)

	result, err := tp.parseType(structType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	// Verify it's a struct type
	if _, ok := result.(*model.TypeStruct); !ok {
		t.Errorf("Expected TypeStruct, got %T", result)
	}
}

func TestParseNamedType(t *testing.T) {
	// Create a simple named type using go/types
	pkg := types.NewPackage("example.com/test", "test")
	underlying := types.NewStruct(nil, nil)
	namedType := types.NewNamed(types.NewTypeName(token.NoPos, pkg, "TestType", nil), underlying, nil)

	p := NewParser()
	tp := p.newTypeParser()

	result, err := tp.parseType(namedType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	// Verify it's a named type
	if namedResult, ok := result.(*model.TypeNamed); ok {
		if namedResult.Name() != "TestType" {
			t.Errorf("Name() = %v, want %v", namedResult.Name(), "TestType")
		}
	} else {
		t.Errorf("Expected TypeNamed, got %T", result)
	}
}

func TestTypeParserCaching(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Parse the same type twice
	stringType := types.Typ[types.String]

	result1, err1 := tp.parseType(stringType)
	if err1 != nil {
		t.Fatalf("First parseType() error = %v", err1)
	}

	result2, err2 := tp.parseType(stringType)
	if err2 != nil {
		t.Fatalf("Second parseType() error = %v", err2)
	}

	// Should return the same instance (cached) or equivalent results
	if result1.PrintType("", model.PackageMap{}) != result2.PrintType("", model.PackageMap{}) {
		t.Error("parseType() should return consistent results for same type")
	}
}

func TestParseArrayType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create array type [5]string
	arrayType := types.NewArray(types.Typ[types.String], 5)

	result, err := tp.parseType(arrayType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	typeStr := result.PrintType("", model.PackageMap{})
	expected := "[5]string"
	if typeStr != expected {
		t.Errorf("PrintType() = %v, want %v", typeStr, expected)
	}
}

func TestParseVariadicType(t *testing.T) {
	p := NewParser()
	tp := p.newTypeParser()

	// Create function type with variadic parameter func(...string)
	params := types.NewTuple(
		types.NewVar(token.NoPos, nil, "args", types.NewSlice(types.Typ[types.String])),
	)
	funcType := types.NewSignatureType(nil, nil, nil, params, nil, true)

	result, err := tp.parseType(funcType)
	if err != nil {
		t.Fatalf("parseType() error = %v", err)
	}

	if result == nil {
		t.Fatal("parseType() returned nil")
	}

	// Check that it's a function type
	if _, ok := result.(*model.TypeSignature); !ok {
		t.Errorf("Expected TypeSignature, got %T", result)
	}
}