package model

import (
	"fmt"
	"strings"
)

// Type is type.
type Type interface {
	PrintType(myPkgPath string, pm PackageMap) string
	addImports(pm *PackageMap)
}

// TypeArray is an array or slice type.
type TypeArray struct {
	Len  int64 // -1 for slices, >= 0 for arrays
	Type Type
}

// PrintType returns type.
func (t *TypeArray) PrintType(myPkgPath string, pm PackageMap) string {
	s := "[]"
	if t.Len > -1 {
		s = fmt.Sprintf("[%d]", t.Len)
	}
	return s + t.Type.PrintType(myPkgPath, pm)
}

func (t *TypeArray) addImports(pm *PackageMap) {
	t.Type.addImports(pm)
}

// TypeBasic represents all types that are not composed from simpler types.
// strings, booleans, and numbers
type TypeBasic string

// PrintType returns type.
func (t *TypeBasic) PrintType(myPkgPath string, pm PackageMap) string {
	return string(*t)
}

func (t *TypeBasic) addImports(pm *PackageMap) {
	// do nothing
}

// TypeChan is map type.
type TypeChan struct {
	Dir  ChanDir
	Type Type
}

// A ChanDir value indicates a channel direction.
type ChanDir int

// The direction of a channel is indicated by one of these constants.
const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)

// PrintType returns type.
func (t *TypeChan) PrintType(myPkgPath string, pm PackageMap) string {
	s := t.Type.PrintType(myPkgPath, pm)
	switch t.Dir {
	case SendRecv:
		return "chan " + s
	case SendOnly:
		return "chan<- " + s
	case RecvOnly:
		return "<-chan " + s
	default:
		return "chan " + s
	}
}

func (t *TypeChan) addImports(pm *PackageMap) {
	t.Type.addImports(pm)
}

// TypeInterface is interface type
type TypeInterface struct {
	Embeddeds       []Type
	ExplicitMethods []*Func
}

// PrintType returns type.
func (t *TypeInterface) PrintType(myPkgPath string, pm PackageMap) string {
	/*
		interface{Embedded ;Func(); Func()}
	*/

	str := "interface{"
	for _, e := range t.Embeddeds {
		str += e.PrintType(myPkgPath, pm)
		str += ";"
	}
	for _, e := range t.ExplicitMethods {
		str += e.PrintDef(myPkgPath, pm)
		str += ";"
	}
	str = strings.TrimRight(str, ";")
	str += "}"

	return str
}

func (t *TypeInterface) addImports(pm *PackageMap) {
	for _, e := range t.Embeddeds {
		e.addImports(pm)
	}
	for _, e := range t.ExplicitMethods {
		e.addImports(pm)
	}
}

// TypeMap is map type.
type TypeMap struct {
	Key   Type
	Value Type
}

// PrintType returns type.
func (t *TypeMap) PrintType(myPkgPath string, pm PackageMap) string {
	return fmt.Sprintf("map[%s]%s", t.Key.PrintType(myPkgPath, pm), t.Value.PrintType(myPkgPath, pm))
}

func (t *TypeMap) addImports(pm *PackageMap) {
	t.Key.addImports(pm)
	t.Value.addImports(pm)
}

// TypeNamed is named.
type TypeNamed struct {
	Pkg  *PkgInfo
	Type string
}

// NewTypeNamed returns TypeNamed.
func NewTypeNamed(pkg *PkgInfo, typ string) Type {
	return &TypeNamed{
		Pkg:  pkg,
		Type: typ,
	}
}

// PrintType returns type.
func (t *TypeNamed) PrintType(myPkgPath string, pm PackageMap) string {
	if t.Pkg == nil {
		return t.Type
	}
	pkg := pm.Get(t.Pkg.Path)
	if pkg == nil {
		// TODO: is internal error?
		pkg = t.Pkg
	}
	return pkg.Prefix(myPkgPath) + t.Type
}

func (t *TypeNamed) addImports(pm *PackageMap) {
	if t.Pkg != nil {
		pm.Need(t.Pkg.Path, true)
	}
}

// TypePointer is a pointer.
type TypePointer struct {
	Type Type
}

// NewPointer returns pointer of the type.
func NewPointer(typ Type) Type {
	return &TypePointer{
		Type: typ,
	}
}

// PrintType returns type.
func (t *TypePointer) PrintType(myPkgPath string, pm PackageMap) string {
	return "*" + t.Type.PrintType(myPkgPath, pm)
}

func (t *TypePointer) addImports(pm *PackageMap) {
	t.Type.addImports(pm)
}

// TypeSignature is function.
type TypeSignature struct {
	Params   []*Parameter
	Variadic *Parameter
	Results  []*Parameter
}

// NewTypeSignature returns TypeSignature.
func NewTypeSignature(params []*Parameter, variadic *Parameter, results []*Parameter) *TypeSignature {
	return &TypeSignature{
		Params:   params,
		Variadic: variadic,
		Results:  results,
	}
}

// PrintType returns type.
func (t *TypeSignature) PrintType(myPkgPath string, pm PackageMap) string {
	s := "func"
	s += t.printArgs(myPkgPath, pm)
	s += t.printResults(myPkgPath, pm)
	return s
}

// PrintParams print params
// for example : (x int, y int)
func (t *TypeSignature) printArgs(myPkgPath string, pm PackageMap) string {
	// args
	s := "("
	for _, param := range t.Params {
		s += param.PrintNameAndType(myPkgPath, pm)
		s += ","
	}
	s = strings.TrimRight(s, ",")
	if t.Variadic != nil {
		s += fmt.Sprintf(",%s ...%s", t.Variadic.Name, t.Variadic.Type.PrintType(myPkgPath, pm))
	}
	s += ")"
	return s
}

func (t *TypeSignature) printResults(myPkgPath string, pm PackageMap) string {
	s := ""
	if len(t.Results) > 1 {
		s += "("
	}
	for i, result := range t.Results {
		if i != 0 {
			s += ","
		}
		s += result.PrintNameAndType(myPkgPath, pm)
	}
	if len(t.Results) > 1 {
		s += ")"
	}
	return s
}

func (t *TypeSignature) addImports(pm *PackageMap) {
	for _, p := range t.Params {
		p.addImports(pm)
	}
	if t.Variadic != nil {
		t.Variadic.addImports(pm)
	}
	for _, p := range t.Results {
		p.addImports(pm)
	}
}

// PrintCallArgsFmt returns format to call this function.
// For exapmle, if the func has tow params, returns "(%s, %s)"
func (t *TypeSignature) PrintCallArgsFmt() string {
	fmt := "("
	for i := 0; i < len(t.Params); i++ {
		fmt += "%s,"
	}
	if t.Variadic != nil {
		fmt += "%s,"
	}
	fmt = strings.TrimRight(fmt, ",")
	fmt += ")"
	return fmt
}

// TypeStruct is struct type.
type TypeStruct struct {
	Fields []*Field
}

// PrintType returns type.
func (t *TypeStruct) PrintType(myPkgPath string, pm PackageMap) string {
	/*
		struct{n XX; n XX}
	*/
	str := "struct{"
	for _, f := range t.Fields {
		str += f.PrintDef(myPkgPath, pm)
		str += ";"
	}
	str = strings.TrimRight(str, ";")
	str += "}"

	return str
}

func (t *TypeStruct) addImports(pm *PackageMap) {
	for _, f := range t.Fields {
		f.addImports(pm)
	}
}
