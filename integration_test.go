package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmio11/codegen/generator/model"
)

// TestIntegrationBasicMockGeneration tests basic mock generation functionality
func TestIntegrationBasicMockGeneration(t *testing.T) {
	// Test creating basic mock structures
	pkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	
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
		model.NewFunc("Save", model.NewTypeSignature(
			[]*model.Parameter{model.NewParameter("data", model.NewTypeBasic("string"))},
			nil,
			[]*model.Parameter{model.NewParameter("", model.NewTypeBasic("error"))},
		), ""),
	}
	
	intf := model.NewInterface("TestInterface", pkg, methods)
	
	// Test mock struct creation
	mockStruct := model.NewStruct("MockTestInterface", pkg)
	if mockStruct == nil {
		t.Fatal("Failed to create mock struct")
	}
	
	if mockStruct.Name() != "MockTestInterface" {
		t.Errorf("Mock struct name = %v, want %v", mockStruct.Name(), "MockTestInterface")
	}
	
	// Test that interface methods are accessible
	if len(intf.Methods()) != 2 {
		t.Errorf("Interface methods count = %v, want %v", len(intf.Methods()), 2)
	}
	
	// Check specific methods
	getMethod := intf.Methods()[0]
	if getMethod.Name() != "Get" {
		t.Errorf("First method name = %v, want %v", getMethod.Name(), "Get")
	}
	
	saveMethod := intf.Methods()[1]
	if saveMethod.Name() != "Save" {
		t.Errorf("Second method name = %v, want %v", saveMethod.Name(), "Save")
	}
}

// TestIntegrationGenericMockStructure tests generic mock structure creation
func TestIntegrationGenericMockStructure(t *testing.T) {
	pkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	// Create type parameters
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
	
	genericIntf := model.NewGenericInterface("Repository", pkg, methods, typeParams)
	
	// Test generic interface properties
	if !genericIntf.IsGeneric() {
		t.Error("Generic interface should be marked as generic")
	}
	
	if len(genericIntf.TypeParams()) != 1 {
		t.Errorf("Generic interface type params count = %v, want %v", len(genericIntf.TypeParams()), 1)
	}
	
	if genericIntf.TypeParams()[0].Name() != "T" {
		t.Errorf("Type parameter name = %v, want %v", genericIntf.TypeParams()[0].Name(), "T")
	}
	
	// Create corresponding generic mock struct
	mockStruct := model.NewGenericStruct("MockRepository", pkg, typeParams)
	
	if !mockStruct.IsGeneric() {
		t.Error("Generic mock struct should be marked as generic")
	}
	
	if len(mockStruct.TypeParams()) != 1 {
		t.Errorf("Generic mock struct type params count = %v, want %v", len(mockStruct.TypeParams()), 1)
	}
}

// TestIntegrationFileGeneration tests complete file generation
func TestIntegrationFileGeneration(t *testing.T) {
	pkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	dependencies := model.NewPackageMap("testpkg", "example.com/testpkg")
	
	// Create a file
	file := model.NewFile("/tmp/test_mock.go", "testpkg", "example.com/testpkg", dependencies)
	
	if file == nil {
		t.Fatal("Failed to create file")
	}
	
	// Create a simple struct to add to the file
	mockStruct := model.NewStruct("MockTestInterface", pkg)
	file.AddStruct(mockStruct)
	
	// Generate code
	code := file.PrintCode()
	
	// Check generated code structure
	if !strings.Contains(code, "package testpkg") {
		t.Error("Generated code should contain package declaration")
	}
	
	if !strings.Contains(code, "type MockTestInterface struct") {
		t.Error("Generated code should contain struct definition")
	}
}

// TestIntegrationGenericFileGeneration tests generic file generation
func TestIntegrationGenericFileGeneration(t *testing.T) {
	pkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	dependencies := model.NewPackageMap("testpkg", "example.com/testpkg")
	
	// Create a file
	file := model.NewFile("/tmp/generic_mock.go", "testpkg", "example.com/testpkg", dependencies)
	
	// Create generic type parameters
	typeParams := []*model.TypeParameter{
		model.NewTypeParameter("T", model.ConstraintAny, 0),
		model.NewTypeParameter("K", model.ConstraintComparable, 1),
	}
	
	// Create a generic struct
	genericStruct := model.NewGenericStruct("MockRepository", pkg, typeParams)
	file.AddStruct(genericStruct)
	
	// Generate code
	code := file.PrintCode()
	
	// Check generated code structure
	if !strings.Contains(code, "package testpkg") {
		t.Error("Generated code should contain package declaration")
	}
	
	if !strings.Contains(code, "type MockRepository[T any, K comparable] struct") {
		t.Error("Generated code should contain generic struct definition")
	}
}

// TestIntegrationMockMethodGeneration tests mock method generation
func TestIntegrationMockMethodGeneration(t *testing.T) {
	// Create mock method signature
	params := []*model.Parameter{
		model.NewParameter("id", model.NewTypeBasic("string")),
	}
	results := []*model.Parameter{
		model.NewParameter("", model.NewTypeBasic("string")),
		model.NewParameter("", model.NewTypeBasic("error")),
	}
	
	signature := model.NewTypeSignature(params, nil, results)
	method := model.NewFunc("Get", signature, "")
	
	// Test method properties
	if method.Name() != "Get" {
		t.Errorf("Method name = %v, want %v", method.Name(), "Get")
	}
	
	if len(method.Type().Args()) != 1 {
		t.Errorf("Method args count = %v, want %v", len(method.Type().Args()), 1)
	}
	
	if len(method.Type().Results()) != 2 {
		t.Errorf("Method results count = %v, want %v", len(method.Type().Results()), 2)
	}
	
	// Check parameter types
	idParam := method.Type().Args()[0]
	if idParam.Name() != "id" {
		t.Errorf("Parameter name = %v, want %v", idParam.Name(), "id")
	}
	
	if idParam.Type().PrintType("", model.PackageMap{}) != "string" {
		t.Errorf("Parameter type = %v, want %v", idParam.Type().PrintType("", model.PackageMap{}), "string")
	}
}

// TestIntegrationComplexTypeGeneration tests complex type generation
func TestIntegrationComplexTypeGeneration(t *testing.T) {
	// Test slice type
	sliceType := model.NewTypeArray(-1, model.NewTypeBasic("string")) // -1 for slice
	sliceStr := sliceType.PrintType("", model.PackageMap{})
	if sliceStr != "[]string" {
		t.Errorf("Slice type = %v, want %v", sliceStr, "[]string")
	}
	
	// Test map type
	mapType := model.NewTypeMap(model.NewTypeBasic("string"), model.NewTypeBasic("int"))
	mapStr := mapType.PrintType("", model.PackageMap{})
	if mapStr != "map[string]int" {
		t.Errorf("Map type = %v, want %v", mapStr, "map[string]int")
	}
	
	// Test pointer type
	ptrType := model.NewPointer(model.NewTypeBasic("string"))
	ptrStr := ptrType.PrintType("", model.PackageMap{})
	if ptrStr != "*string" {
		t.Errorf("Pointer type = %v, want %v", ptrStr, "*string")
	}
	
	// Test channel type
	chanType := model.NewTypeChan(model.SendRecv, model.NewTypeBasic("string"))
	chanStr := chanType.PrintType("", model.PackageMap{})
	if strings.Contains(chanStr, "chan") {
		// Channel type exists
	} else {
		t.Errorf("Channel type should contain 'chan'")
	}
}

// TestIntegrationConstraintTypes tests constraint type generation
func TestIntegrationConstraintTypes(t *testing.T) {
	// Test built-in constraints
	anyConstraint := model.ConstraintAny
	if anyConstraint.Name() != "any" {
		t.Errorf("Any constraint name = %v, want %v", anyConstraint.Name(), "any")
	}
	
	comparableConstraint := model.ConstraintComparable
	if comparableConstraint.Name() != "comparable" {
		t.Errorf("Comparable constraint name = %v, want %v", comparableConstraint.Name(), "comparable")
	}
	
	// Test custom constraint
	customConstraint := model.NewTypeConstraint("Number")
	if customConstraint.Name() != "Number" {
		t.Errorf("Custom constraint name = %v, want %v", customConstraint.Name(), "Number")
	}
	
	customStr := customConstraint.PrintType("", model.PackageMap{})
	if customStr != "Number" {
		t.Errorf("Custom constraint string = %v, want %v", customStr, "Number")
	}
}

// TestIntegrationPackageDependencies tests package dependency management
func TestIntegrationPackageDependencies(t *testing.T) {
	dependencies := model.NewPackageMap("testpkg", "example.com/testpkg")
	
	// Add some dependencies
	dependencies.Add("fmt", *model.NewPkgInfo("fmt", "fmt", ""))
	dependencies.Add("strings", *model.NewPkgInfo("strings", "strings", ""))
	dependencies.Add("context", *model.NewPkgInfo("context", "context", ""))
	
	// Test package retrieval
	fmtPkg := dependencies.Get("fmt")
	if fmtPkg == nil {
		t.Error("Should be able to retrieve added package")
	}
	
	if fmtPkg.Name() != "fmt" {
		t.Errorf("Package name = %v, want %v", fmtPkg.Name(), "fmt")
	}
	
	// Test import generation
	importCode := dependencies.PrintCode("example.com/testpkg")
	if !strings.Contains(importCode, "import (") {
		t.Error("Import code should contain import statement")
	}
	
	if !strings.Contains(importCode, "\"fmt\"") {
		t.Error("Import code should contain fmt import")
	}
}

// TestIntegrationEndToEndWorkflow tests complete workflow
func TestIntegrationEndToEndWorkflow(t *testing.T) {
	// Create temporary directory for test output
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "mock_test.go")
	
	// Create package and dependencies
	pkg := model.NewPkgInfo("testpkg", "example.com/testpkg", "")
	dependencies := model.NewPackageMap("testpkg", "example.com/testpkg")
	dependencies.Add("context", *model.NewPkgInfo("context", "context", ""))
	dependencies.Add("fmt", *model.NewPkgInfo("fmt", "fmt", ""))
	
	// Create file
	file := model.NewFile(outputFile, "testpkg", "example.com/testpkg", dependencies)
	
	// Create interface methods
	methods := []*model.Func{
		model.NewFunc("Process", model.NewTypeSignature(
			[]*model.Parameter{
				model.NewParameter("ctx", model.NewTypeNamed(
					model.NewPkgInfo("context", "context", ""),
					"Context",
					model.NewTypeInterface([]*model.TypeNamed{}, []*model.Func{}),
				)),
				model.NewParameter("data", model.NewTypeBasic("string")),
			},
			nil,
			[]*model.Parameter{model.NewParameter("", model.NewTypeBasic("error"))},
		), ""),
	}
	
	// Create interface
	intf := model.NewInterface("Processor", pkg, methods)
	file.AddInterface(intf)
	
	// Create mock struct
	mockStruct := model.NewStruct("MockProcessor", pkg)
	
	// Add fake method field
	fakeMethodType := model.NewTypeSignature(
		[]*model.Parameter{
			model.NewParameter("ctx", model.NewTypeNamed(
				model.NewPkgInfo("context", "context", ""),
				"Context",
				model.NewTypeInterface([]*model.TypeNamed{}, []*model.Func{}),
			)),
			model.NewParameter("data", model.NewTypeBasic("string")),
		},
		nil,
		[]*model.Parameter{model.NewParameter("", model.NewTypeBasic("error"))},
	)
	
	fakeField := model.NewField("FakeProcess", fakeMethodType, "")
	mockStruct.AddField(fakeField)
	
	// Add mock method implementation  
	mockMethodBody := "return m.FakeProcess(ctx, data)"
	receiverType := model.NewTypeNamed(pkg, "MockProcessor", model.NewTypeStruct([]*model.Field{}))
	receiver := model.NewParameter("m", model.NewPointer(receiverType))
	mockMethod := model.NewMethod(receiver, "Process", model.NewTypeSignature(
		[]*model.Parameter{
			model.NewParameter("ctx", model.NewTypeNamed(
				model.NewPkgInfo("context", "context", ""),
				"Context",
				model.NewTypeInterface([]*model.TypeNamed{}, []*model.Func{}),
			)),
			model.NewParameter("data", model.NewTypeBasic("string")),
		},
		nil,
		[]*model.Parameter{model.NewParameter("", model.NewTypeBasic("error"))},
	), mockMethodBody)
	
	mockStruct.AddMethod(mockMethod)
	file.AddStruct(mockStruct)
	
	// Tidy dependencies
	file.DependenciesTidy()
	
	// Generate code
	code := file.PrintCode()
	
	// Verify generated code
	if !strings.Contains(code, "package testpkg") {
		t.Error("Generated code should contain package declaration")
	}
	
	if !strings.Contains(code, "import (") {
		t.Error("Generated code should contain imports")
	}
	
	if !strings.Contains(code, "\"context\"") {
		t.Error("Generated code should import context package")
	}
	
	if !strings.Contains(code, "type Processor interface") {
		t.Error("Generated code should contain interface definition")
	}
	
	if !strings.Contains(code, "type MockProcessor struct") {
		t.Error("Generated code should contain mock struct")
	}
	
	if !strings.Contains(code, "FakeProcess") {
		t.Error("Generated code should contain fake method field")
	}
	
	// Write to file for verification
	err := os.WriteFile(outputFile, []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file should be created")
	}
}