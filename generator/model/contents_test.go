package model

import (
	"testing"
)

func TestGenericInterface(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	typeParams := []*TypeParameter{
		NewTypeParameter("T", ConstraintAny, 0),
	}
	
	// Create methods for the interface
	methods := []*Func{
		NewFunc("Get", NewTypeSignature(
			[]*Parameter{NewParameter("id", NewTypeBasic("string"))},
			nil,
			[]*Parameter{
				NewParameter("", NewTypeParameter("T", nil, 0)),
				NewParameter("", NewTypeBasic("error")),
			},
		), ""),
		NewFunc("Save", NewTypeSignature(
			[]*Parameter{NewParameter("item", NewTypeParameter("T", nil, 0))},
			nil,
			[]*Parameter{NewParameter("", NewTypeBasic("error"))},
		), ""),
	}
	
	// Test generic interface creation
	genericIntf := NewGenericInterface("Repository", pkg, methods, typeParams)
	
	if !genericIntf.IsGeneric() {
		t.Error("IsGeneric() should return true for generic interface")
	}
	
	if len(genericIntf.TypeParams()) != 1 {
		t.Errorf("TypeParams() length = %v, want %v", len(genericIntf.TypeParams()), 1)
	}
	
	if genericIntf.TypeParams()[0].Name() != "T" {
		t.Errorf("Type parameter name = %v, want %v", genericIntf.TypeParams()[0].Name(), "T")
	}
	
	if genericIntf.Name() != "Repository" {
		t.Errorf("Name() = %v, want %v", genericIntf.Name(), "Repository")
	}
	
	if len(genericIntf.Methods()) != 2 {
		t.Errorf("Methods() length = %v, want %v", len(genericIntf.Methods()), 2)
	}
}

func TestGenericStruct(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	typeParams := []*TypeParameter{
		NewTypeParameter("T", ConstraintAny, 0),
		NewTypeParameter("K", ConstraintComparable, 1),
	}
	
	// Test generic struct creation
	genericStruct := NewGenericStruct("GenericMap", pkg, typeParams)
	
	if !genericStruct.IsGeneric() {
		t.Error("IsGeneric() should return true for generic struct")
	}
	
	if len(genericStruct.TypeParams()) != 2 {
		t.Errorf("TypeParams() length = %v, want %v", len(genericStruct.TypeParams()), 2)
	}
	
	if genericStruct.Name() != "GenericMap" {
		t.Errorf("Name() = %v, want %v", genericStruct.Name(), "GenericMap")
	}
	
	// Add fields with generic types
	genericStruct.AddField(NewField("data", NewTypeMap(
		NewTypeParameter("K", nil, 1),
		NewTypeParameter("T", nil, 0),
	), ""))
	
	if len(genericStruct.Fields()) != 1 {
		t.Errorf("Fields() length = %v, want %v", len(genericStruct.Fields()), 1)
	}
}

func TestNonGenericInterface(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	// Create methods for non-generic interface
	methods := []*Func{
		NewFunc("Process", NewTypeSignature(
			[]*Parameter{NewParameter("data", NewTypeBasic("string"))},
			nil,
			[]*Parameter{NewParameter("", NewTypeBasic("error"))},
		), ""),
	}
	
	// Test non-generic interface creation
	nonGenericIntf := NewInterface("Processor", pkg, methods)
	
	if nonGenericIntf.IsGeneric() {
		t.Error("IsGeneric() should return false for non-generic interface")
	}
	
	if len(nonGenericIntf.TypeParams()) != 0 {
		t.Errorf("TypeParams() length = %v, want %v", len(nonGenericIntf.TypeParams()), 0)
	}
	
	if nonGenericIntf.Name() != "Processor" {
		t.Errorf("Name() = %v, want %v", nonGenericIntf.Name(), "Processor")
	}
}

func TestNonGenericStruct(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	// Test non-generic struct creation
	nonGenericStruct := NewStruct("SimpleStruct", pkg)
	
	if nonGenericStruct.IsGeneric() {
		t.Error("IsGeneric() should return false for non-generic struct")
	}
	
	if len(nonGenericStruct.TypeParams()) != 0 {
		t.Errorf("TypeParams() length = %v, want %v", len(nonGenericStruct.TypeParams()), 0)
	}
	
	if nonGenericStruct.Name() != "SimpleStruct" {
		t.Errorf("Name() = %v, want %v", nonGenericStruct.Name(), "SimpleStruct")
	}
}