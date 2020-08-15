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
	len int64 // -1 for slices, >= 0 for arrays
	typ Type
}

// NewTypeArray returns TypeArray.
func NewTypeArray(len int64, typ Type) *TypeArray {
	return &TypeArray{
		len: len,
		typ: typ,
	}
}

// Len returns len
func (t *TypeArray) Len() int64 {
	return t.len
}

// Type returns Type.
func (t *TypeArray) Type() Type {
	return t.typ
}

// PrintType returns type.
func (t *TypeArray) PrintType(myPkgPath string, pm PackageMap) string {
	s := "[]"
	if t.Len() > -1 {
		s = fmt.Sprintf("[%d]", t.Len())
	}
	return s + t.Type().PrintType(myPkgPath, pm)
}

func (t *TypeArray) addImports(pm *PackageMap) {
	t.Type().addImports(pm)
}

// TypeBasic represents all types that are not composed from simpler types.
// strings, booleans, and numbers
type TypeBasic string

// NewTypeBasic returns TypeBasic.
func NewTypeBasic(s string) *TypeBasic {
	ss := TypeBasic(s)
	return &ss
}

// PrintType returns type.
func (t *TypeBasic) PrintType(myPkgPath string, pm PackageMap) string {
	return string(*t)
}

func (t *TypeBasic) addImports(pm *PackageMap) {
	// do nothing
}

// TypeChan is map type.
type TypeChan struct {
	dir ChanDir
	typ Type
}

// A ChanDir value indicates a channel direction.
type ChanDir int

// The direction of a channel is indicated by one of these constants.
const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)

// NewTypeChan returns TypeChan.
func NewTypeChan(dir ChanDir, typ Type) *TypeChan {
	return &TypeChan{
		dir: dir,
		typ: typ,
	}
}

// Dir returns a channel direction.
func (t *TypeChan) Dir() ChanDir {
	return t.dir
}

// Type returns type.
func (t *TypeChan) Type() Type {
	return t.typ
}

// PrintType returns type.
func (t *TypeChan) PrintType(myPkgPath string, pm PackageMap) string {
	s := t.typ.PrintType(myPkgPath, pm)
	switch t.dir {
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
	t.typ.addImports(pm)
}

// TypeInterface is interface type
type TypeInterface struct {
	embeddeds       []Type
	explicitMethods []*Func
}

// NewTypeInterface returns TypeInterface.
func NewTypeInterface(embeddeds []Type, exMethods []*Func) *TypeInterface {
	return &TypeInterface{
		embeddeds:       embeddeds,
		explicitMethods: exMethods,
	}
}

// Embeddeds returns Embedded types.
func (t *TypeInterface) Embeddeds() []Type {
	return t.embeddeds
}

// ExplicitMethods returns explicit methods.
func (t *TypeInterface) ExplicitMethods() []*Func {
	return t.explicitMethods
}

// PrintType returns type.
func (t *TypeInterface) PrintType(myPkgPath string, pm PackageMap) string {
	/*
		interface{Embedded ;Func(); Func()}
	*/

	str := "interface{"
	for _, e := range t.embeddeds {
		str += e.PrintType(myPkgPath, pm)
		str += ";"
	}
	for _, e := range t.explicitMethods {
		str += e.PrintDef(myPkgPath, pm)
		str += ";"
	}
	str = strings.TrimRight(str, ";")
	str += "}"

	return str
}

func (t *TypeInterface) addImports(pm *PackageMap) {
	for _, e := range t.embeddeds {
		e.addImports(pm)
	}
	for _, e := range t.explicitMethods {
		e.addImports(pm)
	}
}

// TypeMap is map type.
type TypeMap struct {
	key   Type
	value Type
}

// NewTypeMap returns TypeMap.
func NewTypeMap(key Type, val Type) *TypeMap {
	return &TypeMap{
		key:   key,
		value: val,
	}
}

// Key returns key type.
func (t *TypeMap) Key() Type {
	return t.key
}

// Value returns value type.
func (t *TypeMap) Value() Type {
	return t.value
}

// PrintType returns type.
func (t *TypeMap) PrintType(myPkgPath string, pm PackageMap) string {
	return fmt.Sprintf("map[%s]%s", t.key.PrintType(myPkgPath, pm), t.value.PrintType(myPkgPath, pm))
}

func (t *TypeMap) addImports(pm *PackageMap) {
	t.key.addImports(pm)
	t.value.addImports(pm)
}

// TypeNamed is named.
type TypeNamed struct {
	pkg  *PkgInfo
	name string
}

// NewTypeNamed returns TypeNamed.
func NewTypeNamed(pkg *PkgInfo, name string) Type {
	return &TypeNamed{
		pkg:  pkg,
		name: name,
	}
}

// Pkg returns package info.
func (t *TypeNamed) Pkg() *PkgInfo {
	return t.pkg
}

// Name returns name.
func (t *TypeNamed) Name() string {
	return t.name
}

// PrintType returns type.
func (t *TypeNamed) PrintType(myPkgPath string, pm PackageMap) string {
	if t.pkg == nil {
		return t.name
	}
	pkg := pm.Get(t.pkg.Path)
	if pkg == nil {
		// TODO: is internal error?
		pkg = t.pkg
	}
	return pkg.Prefix(myPkgPath) + t.name
}

func (t *TypeNamed) addImports(pm *PackageMap) {
	if t.pkg != nil {
		pm.Need(t.pkg.Path, true)
	}
}

// TypePointer is a pointer.
type TypePointer struct {
	typ Type
}

// NewPointer returns pointer of the type.
func NewPointer(typ Type) Type {
	return &TypePointer{
		typ: typ,
	}
}

// Type returns type.
func (t *TypePointer) Type() Type {
	return t.typ
}

// PrintType returns type.
func (t *TypePointer) PrintType(myPkgPath string, pm PackageMap) string {
	return "*" + t.typ.PrintType(myPkgPath, pm)
}

func (t *TypePointer) addImports(pm *PackageMap) {
	t.typ.addImports(pm)
}

// TypeSignature is function.
type TypeSignature struct {
	args     []*Parameter
	variadic *Parameter
	results  []*Parameter
}

// NewTypeSignature returns TypeSignature.
func NewTypeSignature(params []*Parameter, variadic *Parameter, results []*Parameter) *TypeSignature {
	return &TypeSignature{
		args:     params,
		variadic: variadic,
		results:  results,
	}
}

// Args returns arguments.
func (t *TypeSignature) Args() []*Parameter {
	return t.args
}

// Variadic returns variadic parameter.
// If nil, the function has no variadic parameter.
func (t *TypeSignature) Variadic() *Parameter {
	return t.variadic
}

// Results returns results.
func (t *TypeSignature) Results() []*Parameter {
	return t.results
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
	for _, param := range t.args {
		s += param.PrintNameAndType(myPkgPath, pm)
		s += ","
	}
	s = strings.TrimRight(s, ",")
	if t.variadic != nil {
		s += fmt.Sprintf(",%s ...%s", t.variadic.Name(), t.variadic.Type().PrintType(myPkgPath, pm))
	}
	s += ")"
	return s
}

func (t *TypeSignature) printResults(myPkgPath string, pm PackageMap) string {
	s := ""
	if len(t.results) > 1 {
		s += "("
	}
	for i, result := range t.results {
		if i != 0 {
			s += ","
		}
		s += result.PrintNameAndType(myPkgPath, pm)
	}
	if len(t.results) > 1 {
		s += ")"
	}
	return s
}

func (t *TypeSignature) addImports(pm *PackageMap) {
	for _, p := range t.args {
		p.addImports(pm)
	}
	if t.variadic != nil {
		t.variadic.addImports(pm)
	}
	for _, p := range t.results {
		p.addImports(pm)
	}
}

// PrintCallArgsFmt returns format to call this function.
// For exapmle, if the func has two arguments, returns "(%s, %s)"
func (t *TypeSignature) PrintCallArgsFmt() string {
	fmt := "("
	for i := 0; i < len(t.args); i++ {
		fmt += "%s,"
	}
	if t.variadic != nil {
		fmt += "%s,"
	}
	fmt = strings.TrimRight(fmt, ",")
	fmt += ")"
	return fmt
}

// TypeStruct is struct type.
type TypeStruct struct {
	fields []*Field
}

// NewTypeStruct returns TypeStruct.
func NewTypeStruct(fields []*Field) *TypeStruct {
	return &TypeStruct{
		fields: fields,
	}
}

// Fields returns fields.
func (t *TypeStruct) Fields() []*Field {
	return t.fields
}

// PrintType returns type.
func (t *TypeStruct) PrintType(myPkgPath string, pm PackageMap) string {
	/*
		struct{n XX; n XX}
	*/
	str := "struct{"
	for _, f := range t.fields {
		str += f.PrintDef(myPkgPath, pm)
		str += ";"
	}
	str = strings.TrimRight(str, ";")
	str += "}"

	return str
}

func (t *TypeStruct) addImports(pm *PackageMap) {
	for _, f := range t.fields {
		f.addImports(pm)
	}
}
