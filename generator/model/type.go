package model

import (
	"fmt"
	"strings"
)

// Type is type.
type Type interface {
	// PrintType returns type.
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

// TypeBasic represents all types that are not composed from simpler types.
// strings, booleans, and numbers
type TypeBasic string

// NewTypeBasic returns TypeBasic.
func NewTypeBasic(s string) *TypeBasic {
	ss := TypeBasic(s)
	return &ss
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

// TypeInterface is interface type
type TypeInterface struct {
	embeddeds       []*TypeNamed
	explicitMethods []*Func
}

// NewTypeInterface returns TypeInterface.
func NewTypeInterface(embeddeds []*TypeNamed, exMethods []*Func) *TypeInterface {
	return &TypeInterface{
		embeddeds:       embeddeds,
		explicitMethods: exMethods,
	}
}

// Embeddeds returns Embedded types.
func (t *TypeInterface) Embeddeds() []*TypeNamed {
	return t.embeddeds
}

// ExplicitMethods returns explicit methods.
func (t *TypeInterface) ExplicitMethods() []*Func {
	return t.explicitMethods
}

// Methods returns all methods interface has.
func (t *TypeInterface) Methods() []*Func {
	return getIntfMethod(t)
}

// getIntfMethod returns all methods interface has.
func getIntfMethod(t Type) []*Func {
	methods := []*Func{}
	var get func(Type)
	get = func(typ Type) {
		if typ == nil {
			return
		}
		switch tt := typ.(type) {
		case *TypeNamed:
			get(tt.org)
			return
		case *TypeInterface:
			for _, em := range tt.embeddeds {
				get(em)
			}
			methods = append(methods, tt.explicitMethods...)
		default:
			panic("unexpected type.")
		}
	}
	get(t)
	return methods
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

// TypeNamed is named.
type TypeNamed struct {
	pkg  *PkgInfo
	name string
	org  Type
}

// NewTypeNamed returns TypeNamed.
func NewTypeNamed(pkg *PkgInfo, name string, org Type) *TypeNamed {
	return &TypeNamed{
		pkg:  pkg,
		name: name,
		org:  org,
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

// Org returns original type.
func (t *TypeNamed) Org() Type {
	return t.org
}

// PrintTypeDef returns type definition.
func (t *TypeNamed) PrintTypeDef(myPkgPath string, pm PackageMap) string {
	/*
		e.g. type XXX struct{}
	*/
	return fmt.Sprintf("type %s %s", t.name, t.org.PrintType(myPkgPath, pm))
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

// AddField add field to struct
func (t *TypeStruct) AddField(f *Field) {
	t.fields = append(t.fields, f)
}

/*
	Implementations for Type methods.
*/

func (t *TypeArray) addImports(pm *PackageMap)     { typeImports(t, pm) }
func (t *TypeBasic) addImports(pm *PackageMap)     { typeImports(t, pm) }
func (t *TypeChan) addImports(pm *PackageMap)      { typeImports(t, pm) }
func (t *TypeInterface) addImports(pm *PackageMap) { typeImports(t, pm) }
func (t *TypeMap) addImports(pm *PackageMap)       { typeImports(t, pm) }
func (t *TypeNamed) addImports(pm *PackageMap)     { typeImports(t, pm) }
func (t *TypePointer) addImports(pm *PackageMap)   { typeImports(t, pm) }
func (t *TypeSignature) addImports(pm *PackageMap) { typeImports(t, pm) }
func (t *TypeStruct) addImports(pm *PackageMap)    { typeImports(t, pm) }

// PrintType returns type.
func (t *TypeArray) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeBasic) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeChan) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeInterface) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeMap) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeNamed) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypePointer) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeSignature) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}

// PrintType returns type.
func (t *TypeStruct) PrintType(myPkgPath string, pm PackageMap) string {
	return printType(t, myPkgPath, pm)
}
