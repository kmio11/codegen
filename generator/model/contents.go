package model

import (
	"fmt"
)

// Parameter is an argument or return parameter of a method.
type Parameter struct {
	name string
	typ  Type
}

func (p *Parameter) addImports(pm *PackageMap) {
	p.typ.addImports(pm)
}

// NewParameter returns Parameter.
func NewParameter(name string, typ Type) *Parameter {
	return &Parameter{
		name: name,
		typ:  typ,
	}
}

// Name returns name.
func (p *Parameter) Name() string {
	return p.name
}

// Type returns type.
func (p *Parameter) Type() Type {
	return p.typ
}

// PrintNameAndType print code.
// x int
func (p *Parameter) PrintNameAndType(myPkgPath string, pm PackageMap) string {
	return p.name + " " + p.typ.PrintType(myPkgPath, pm)
}

// Contents is contents of source code.
type Contents interface {
	PrintCode(myPkgPath string, pm PackageMap) string
	addImports(pm *PackageMap)
}

// Interface is interface.
type Interface struct {
	typ *TypeNamed
}

// NewInterface returns Interface.
func NewInterface(name string, pkg *PkgInfo, methods []*Func, embeddeds ...*TypeNamed) *Interface {
	return &Interface{
		typ: NewTypeNamed(pkg, name,
			NewTypeInterface(
				embeddeds,
				methods,
			),
		),
	}
}

// NewGenericInterface returns generic Interface with type parameters.
func NewGenericInterface(name string, pkg *PkgInfo, methods []*Func, typeParams []*TypeParameter, embeddeds ...*TypeNamed) *Interface {
	return &Interface{
		typ: NewGenericTypeNamed(pkg, name,
			NewGenericTypeInterface(
				embeddeds,
				methods,
				typeParams,
			),
			typeParams,
		),
	}
}

// Name returns name.
func (i *Interface) Name() string {
	return i.typ.Name()
}

// Type returns type.
func (i *Interface) Type() *TypeNamed {
	return i.typ
}

// Methods returns methods.
func (i Interface) Methods() []*Func {
	return getIntfMethod(i.typ)
}

// TypeParams returns type parameters.
func (i *Interface) TypeParams() []*TypeParameter {
	return i.typ.TypeParams()
}

// IsGeneric returns true if this interface has type parameters.
func (i *Interface) IsGeneric() bool {
	return i.typ.IsGeneric()
}

func (i Interface) addImports(pm *PackageMap) {
	i.typ.addImports(pm)
	i.typ.Org().addImports(pm)
}

// PrintCode print code.
func (i *Interface) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		type Foo interface{
			GetXXX (x int, y int) int
		}
	*/
	return i.typ.PrintTypeDef(myPkgPath, pm)
}

// Struct is struct.
type Struct struct {
	typ     *TypeNamed
	methods []*Method
}

// NewStruct return Struct
func NewStruct(name string, pkg *PkgInfo) *Struct {
	return &Struct{
		typ:     NewTypeNamed(pkg, name, NewTypeStruct([]*Field{})),
		methods: []*Method{},
	}
}

// NewGenericStruct returns generic Struct with type parameters.
func NewGenericStruct(name string, pkg *PkgInfo, typeParams []*TypeParameter) *Struct {
	return &Struct{
		typ:     NewGenericTypeNamed(pkg, name, NewTypeStruct([]*Field{}), typeParams),
		methods: []*Method{},
	}
}

// Name returns name.
func (s *Struct) Name() string {
	return s.typ.Name()
}

// Type returns named type.
func (s *Struct) Type() *TypeNamed {
	return s.typ
}

// TypeStruct returns struct type.
func (s *Struct) TypeStruct() *TypeStruct {
	return s.typ.Org().(*TypeStruct)
}

// Methods returns methods.
func (s *Struct) Methods() []*Method {
	return s.methods
}

// Fields returns fields.
func (s *Struct) Fields() []*Field {
	return s.TypeStruct().Fields()
}

// TypeParams returns type parameters.
func (s *Struct) TypeParams() []*TypeParameter {
	return s.typ.TypeParams()
}

// IsGeneric returns true if this struct has type parameters.
func (s *Struct) IsGeneric() bool {
	return s.typ.IsGeneric()
}

func (s *Struct) addImports(pm *PackageMap) {
	s.typ.addImports(pm)
	s.TypeStruct().addImports(pm)

	for _, m := range s.methods {
		m.addImports(pm)
	}
}

// PrintCode print code.
func (s *Struct) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		type Foo struct{
			x int `xxx`
		}

		func (f *Foo) Baa() int{
			return 1
		}
	*/
	str := s.typ.PrintTypeDef(myPkgPath, pm)

	// methods
	str += "\n"
	for _, m := range s.methods {
		str += "\n"
		str += m.PrintCode(myPkgPath, pm)
	}
	return str
}

// AddField add field to struct
func (s *Struct) AddField(f *Field) {
	s.TypeStruct().AddField(f)
}

// AddMethod add method to struct
func (s *Struct) AddMethod(m *Method) {
	s.methods = append(s.methods, m)
}

// Field is field of struct.
// If Parameter.Name is brank, it represents embeded field.
type Field struct {
	Parameter
	tag string
}

// NewField returns Field
func NewField(name string, typ Type, tag string) *Field {
	p := NewParameter(name, typ)
	return &Field{
		Parameter: Parameter{
			name: p.name,
			typ:  p.typ,
		},
		tag: tag,
	}
}

// Tag returns tag.
func (f *Field) Tag() string {
	return f.tag
}

// PrintDef print name ,type and tag(if exist) .
func (f *Field) PrintDef(myPkgPath string, pm PackageMap) string {
	if f.tag != "" {
		return fmt.Sprintf("%s `%s`", f.Parameter.PrintNameAndType(myPkgPath, pm), f.tag)
	}
	return f.Parameter.PrintNameAndType(myPkgPath, pm)
}

// Method is method.
type Method struct {
	rcv *Parameter
	Func
}

// NewMethod returns Method.
func NewMethod(rcv *Parameter, name string, typ *TypeSignature, statements string) *Method {
	f := NewFunc(name, typ, statements)
	return &Method{
		rcv: rcv,
		Func: Func{
			name:       f.name,
			typ:        f.typ,
			statements: f.statements,
		},
	}
}

// Reciever returns reciever.
func (m *Method) Reciever() *Parameter {
	return m.rcv
}

// PrintCode print code.
func (m *Method) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		func (m Reciever) SomeFunc (x int, y int) (int error){
			return x+i , nil
		}
	*/
	s := "func "
	s += fmt.Sprintf("(%s)", m.rcv.PrintNameAndType(myPkgPath, pm))
	s += m.name
	s += m.typ.printArgs(myPkgPath, pm)
	s += m.typ.printResults(myPkgPath, pm)
	s += "{\n"

	s += m.statements
	s += "\n"

	s += "}\n"
	return s
}

// Func is func.
type Func struct {
	name       string
	typ        *TypeSignature
	statements string
}

// NewFunc returns Func.
func NewFunc(name string, typ *TypeSignature, statements string) *Func {
	return &Func{
		name:       name,
		typ:        typ,
		statements: statements,
	}
}

// Name returns name.
func (f *Func) Name() string {
	return f.name
}

// Type returns type.
func (f *Func) Type() *TypeSignature {
	return f.typ
}

// Statements returns statements.
func (f *Func) Statements() string {
	return f.statements
}

// PrintDef print Name and Params and Results
func (f *Func) PrintDef(myPkgPath string, pm PackageMap) string {
	/*
		Func(x int) int
	*/
	return f.name + f.typ.printArgs(myPkgPath, pm) + f.typ.printResults(myPkgPath, pm)
}

// PrintCode print code.
func (f *Func) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		func SomeFunc (x int, y int) (int error){
			return x+i , nil
		}
	*/
	s := "func "
	s += f.name
	s += f.typ.printArgs(myPkgPath, pm)
	s += f.typ.printResults(myPkgPath, pm)
	s += "{\n"

	s += f.statements
	s += "\n"

	s += "}\n"
	return s
}

func (f *Func) addImports(pm *PackageMap) {
	f.typ.addImports(pm)
}

// SetStatements set statements.
func (f *Func) SetStatements(statements string) {
	f.statements = statements
}
