package mock

import (
	"testing"

	"github.com/kmio11/codegen/generator/model"
)

func TestCommand(t *testing.T) {
	cmd := New()
	
	if cmd == nil {
		t.Fatal("New() returned nil")
	}
	
	if cmd.Name() != "mock" {
		t.Errorf("Name() = %v, want %v", cmd.Name(), "mock")
	}
	
	if cmd.Description() != "generate mock" {
		t.Errorf("Description() = %v, want %v", cmd.Description(), "generate mock")
	}
}

func TestCommandParse(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "valid args",
			args:      []string{"-pkg", ".", "-type", "TestInterface"},
			expectErr: false,
		},
		{
			name:      "missing pkg",
			args:      []string{"-type", "TestInterface"},
			expectErr: false, // -pkg has default value "."
		},
		{
			name:      "missing type",
			args:      []string{"-pkg", "."},
			expectErr: true,
		},
		{
			name:      "empty args",
			args:      []string{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			err := cmd.Parse(tt.args)
			
			if tt.expectErr && err == nil {
				t.Error("Parse() should return error for invalid args")
			}
			
			if !tt.expectErr && err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
			}
		})
	}
}

func TestGetMockFieldName(t *testing.T) {
	tests := []struct {
		methodName string
		expected   string
	}{
		{"Get", "FakeGet"},
		{"Save", "FakeSave"},
		{"Process", "FakeProcess"},
		{"", "Fake"},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			result := getMockFieldName(tt.methodName)
			if result != tt.expected {
				t.Errorf("getMockFieldName(%v) = %v, want %v", tt.methodName, result, tt.expected)
			}
		})
	}
}

func TestGetMockArgsName(t *testing.T) {
	tests := []struct {
		index    int
		expected string
	}{
		{0, "a0"},
		{1, "a1"},
		{10, "a10"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getMockArgsName(tt.index)
			if result != tt.expected {
				t.Errorf("getMockArgsName(%v) = %v, want %v", tt.index, result, tt.expected)
			}
		})
	}
}

func TestGetStubMethodName(t *testing.T) {
	tests := []struct {
		methodName string
		expected   string
	}{
		{"Get", "FakeGet"},
		{"Save", "FakeSave"},
		{"Process", "FakeProcess"},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			result := getStubMethodName(tt.methodName)
			if result != tt.expected {
				t.Errorf("getStubMethodName(%v) = %v, want %v", tt.methodName, result, tt.expected)
			}
		})
	}
}

func TestFmtSignature(t *testing.T) {
	// Create a simple signature for testing
	params := []*model.Parameter{
		model.NewParameter("x", model.NewTypeBasic("int")),
		model.NewParameter("y", model.NewTypeBasic("string")),
	}
	results := []*model.Parameter{
		model.NewParameter("result", model.NewTypeBasic("bool")),
		model.NewParameter("err", model.NewTypeBasic("error")),
	}
	
	originalSig := model.NewTypeSignature(params, nil, results)
	
	formattedSig := fmtSignature(originalSig)
	
	if len(formattedSig.Args()) != 2 {
		t.Errorf("fmtSignature() args count = %v, want %v", len(formattedSig.Args()), 2)
	}
	
	if len(formattedSig.Results()) != 2 {
		t.Errorf("fmtSignature() results count = %v, want %v", len(formattedSig.Results()), 2)
	}
	
	// Check that parameter names were changed to standardized format
	if formattedSig.Args()[0].Name() != "a0" {
		t.Errorf("fmtSignature() first arg name = %v, want %v", formattedSig.Args()[0].Name(), "a0")
	}
	
	if formattedSig.Args()[1].Name() != "a1" {
		t.Errorf("fmtSignature() second arg name = %v, want %v", formattedSig.Args()[1].Name(), "a1")
	}
}

func TestFmtSignatureWithVariadic(t *testing.T) {
	// Create a signature with variadic parameter
	params := []*model.Parameter{
		model.NewParameter("prefix", model.NewTypeBasic("string")),
	}
	variadic := model.NewParameter("items", model.NewTypeBasic("string"))
	results := []*model.Parameter{
		model.NewParameter("", model.NewTypeBasic("string")),
	}
	
	originalSig := model.NewTypeSignature(params, variadic, results)
	
	formattedSig := fmtSignature(originalSig)
	
	if len(formattedSig.Args()) != 1 {
		t.Errorf("fmtSignature() args count = %v, want %v", len(formattedSig.Args()), 1)
	}
	
	if formattedSig.Variadic() == nil {
		t.Error("fmtSignature() should preserve variadic parameter")
	}
	
	if formattedSig.Variadic().Name() != "a1" {
		t.Errorf("fmtSignature() variadic name = %v, want %v", formattedSig.Variadic().Name(), "a1")
	}
}

func TestMockImplBasic(t *testing.T) {
	// Create a simple package and interface for testing
	pkg := &model.Package{
		Name: "testpkg",
		Path: "example.com/testpkg",
	}
	
	outPkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	// Create a simple interface
	methods := []*model.Func{
		model.NewFunc("Get", model.NewTypeSignature(
			[]*model.Parameter{model.NewParameter("id", model.NewTypeBasic("string"))},
			nil,
			[]*model.Parameter{
				model.NewParameter("", model.NewTypeBasic("string")),
				model.NewParameter("", model.NewTypeBasic("error")),
			},
		), ""),
	}
	
	intf := model.NewInterface("TestInterface", outPkg, methods)
	
	// Test mock implementation generation
	mockStruct := mockImpl(pkg, intf, outPkg)
	
	if mockStruct == nil {
		t.Fatal("mockImpl() returned nil")
	}
	
	if mockStruct.Name() != "MockTestInterface" {
		t.Errorf("mockImpl() name = %v, want %v", mockStruct.Name(), "MockTestInterface")
	}
	
	if !mockStruct.IsGeneric() && len(mockStruct.Fields()) < 2 {
		t.Error("mockImpl() should create interface field and fake method fields")
	}
	
	if len(mockStruct.Methods()) != 1 {
		t.Errorf("mockImpl() methods count = %v, want %v", len(mockStruct.Methods()), 1)
	}
}

func TestMockImplGeneric(t *testing.T) {
	// Create a simple package and generic interface for testing
	pkg := &model.Package{
		Name: "testpkg",
		Path: "example.com/testpkg",
	}
	
	outPkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	typeParams := []*model.TypeParameter{
		model.NewTypeParameter("T", model.ConstraintAny, 0),
	}
	
	// Create a generic interface
	methods := []*model.Func{
		model.NewFunc("Get", model.NewTypeSignature(
			[]*model.Parameter{model.NewParameter("id", model.NewTypeBasic("string"))},
			nil,
			[]*model.Parameter{
				model.NewParameter("", model.NewTypeParameter("T", nil, 0)),
				model.NewParameter("", model.NewTypeBasic("error")),
			},
		), ""),
	}
	
	intf := model.NewGenericInterface("Repository", outPkg, methods, typeParams)
	
	// Test generic mock implementation generation
	mockStruct := mockImpl(pkg, intf, outPkg)
	
	if mockStruct == nil {
		t.Fatal("mockImpl() returned nil")
	}
	
	if mockStruct.Name() != "MockRepository" {
		t.Errorf("mockImpl() name = %v, want %v", mockStruct.Name(), "MockRepository")
	}
	
	if !mockStruct.IsGeneric() {
		t.Error("mockImpl() should create generic mock for generic interface")
	}
	
	if len(mockStruct.TypeParams()) != 1 {
		t.Errorf("mockImpl() type params count = %v, want %v", len(mockStruct.TypeParams()), 1)
	}
}