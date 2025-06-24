package model

import (
	"testing"
)

func TestTypeParameter(t *testing.T) {
	tests := []struct {
		name        string
		paramName   string
		constraint  Type
		index       int
		expectedStr string
	}{
		{
			name:        "basic type parameter",
			paramName:   "T",
			constraint:  ConstraintAny,
			index:       0,
			expectedStr: "T",
		},
		{
			name:        "comparable constraint",
			paramName:   "K",
			constraint:  ConstraintComparable,
			index:       1,
			expectedStr: "K",
		},
		{
			name:        "no constraint",
			paramName:   "V",
			constraint:  nil,
			index:       2,
			expectedStr: "V",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := NewTypeParameter(tt.paramName, tt.constraint, tt.index)
			
			if param.Name() != tt.paramName {
				t.Errorf("Name() = %v, want %v", param.Name(), tt.paramName)
			}
			
			if param.Constraint() != tt.constraint {
				t.Errorf("Constraint() = %v, want %v", param.Constraint(), tt.constraint)
			}
			
			if param.Index() != tt.index {
				t.Errorf("Index() = %v, want %v", param.Index(), tt.index)
			}
			
			result := param.PrintType("", PackageMap{})
			if result != tt.expectedStr {
				t.Errorf("PrintType() = %v, want %v", result, tt.expectedStr)
			}
		})
	}
}

func TestTypeConstraint(t *testing.T) {
	tests := []struct {
		name         string
		constraintName string
		expected     string
	}{
		{
			name:         "any constraint",
			constraintName: "any",
			expected:     "any",
		},
		{
			name:         "comparable constraint",
			constraintName: "comparable",
			expected:     "comparable",
		},
		{
			name:         "union constraint",
			constraintName: "int | float64",
			expected:     "int | float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := NewTypeConstraint(tt.constraintName)
			
			if constraint.Name() != tt.constraintName {
				t.Errorf("Name() = %v, want %v", constraint.Name(), tt.constraintName)
			}
			
			result := constraint.PrintType("", PackageMap{})
			if result != tt.expected {
				t.Errorf("PrintType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenericTypeNamed(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	typeParams := []*TypeParameter{
		NewTypeParameter("T", ConstraintAny, 0),
		NewTypeParameter("K", ConstraintComparable, 1),
	}
	
	// Test generic type creation
	genericType := NewGenericTypeNamed(pkg, "GenericType", NewTypeStruct([]*Field{}), typeParams)
	
	if !genericType.IsGeneric() {
		t.Error("IsGeneric() should return true for generic type")
	}
	
	if len(genericType.TypeParams()) != 2 {
		t.Errorf("TypeParams() length = %v, want %v", len(genericType.TypeParams()), 2)
	}
	
	if genericType.TypeParams()[0].Name() != "T" {
		t.Errorf("First type parameter name = %v, want %v", genericType.TypeParams()[0].Name(), "T")
	}
	
	if genericType.TypeParams()[1].Name() != "K" {
		t.Errorf("Second type parameter name = %v, want %v", genericType.TypeParams()[1].Name(), "K")
	}
	
	// Test type definition printing
	typeDef := genericType.PrintTypeDef("", PackageMap{})
	expected := "type GenericType[T any, K comparable] struct {\n}"
	if typeDef != expected {
		t.Errorf("PrintTypeDef() = %v, want %v", typeDef, expected)
	}
}

func TestGenericTypeInterface(t *testing.T) {
	typeParams := []*TypeParameter{
		NewTypeParameter("T", ConstraintAny, 0),
	}
	
	// Create a simple method for the interface
	method := NewFunc("Get", NewTypeSignature(
		[]*Parameter{NewParameter("id", NewTypeBasic("string"))},
		nil,
		[]*Parameter{NewParameter("", NewTypeParameter("T", nil, 0))},
	), "")
	
	// Test generic interface creation
	genericInterface := NewGenericTypeInterface([]*TypeNamed{}, []*Func{method}, typeParams)
	
	if !genericInterface.IsGeneric() {
		t.Error("IsGeneric() should return true for generic interface")
	}
	
	if len(genericInterface.TypeParams()) != 1 {
		t.Errorf("TypeParams() length = %v, want %v", len(genericInterface.TypeParams()), 1)
	}
	
	if genericInterface.TypeParams()[0].Name() != "T" {
		t.Errorf("Type parameter name = %v, want %v", genericInterface.TypeParams()[0].Name(), "T")
	}
}

func TestNonGenericTypes(t *testing.T) {
	pkg := NewPkgInfo("testpkg", "example.com/testpkg", "")
	
	// Test non-generic type
	nonGenericType := NewTypeNamed(pkg, "SimpleType", NewTypeStruct([]*Field{}))
	
	if nonGenericType.IsGeneric() {
		t.Error("IsGeneric() should return false for non-generic type")
	}
	
	if len(nonGenericType.TypeParams()) != 0 {
		t.Errorf("TypeParams() length = %v, want %v", len(nonGenericType.TypeParams()), 0)
	}
	
	// Test type definition printing
	typeDef := nonGenericType.PrintTypeDef("", PackageMap{})
	expected := "type SimpleType struct {\n}"
	if typeDef != expected {
		t.Errorf("PrintTypeDef() = %v, want %v", typeDef, expected)
	}
}

func TestBuiltinConstraints(t *testing.T) {
	if ConstraintAny.Name() != "any" {
		t.Errorf("ConstraintAny.Name() = %v, want %v", ConstraintAny.Name(), "any")
	}
	
	if ConstraintComparable.Name() != "comparable" {
		t.Errorf("ConstraintComparable.Name() = %v, want %v", ConstraintComparable.Name(), "comparable")
	}
}