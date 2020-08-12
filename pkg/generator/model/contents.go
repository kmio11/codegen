package model

import (
	"fmt"
)

// Parameter is an argument or return parameter of a method.
type Parameter struct {
	Name string
	Type Type
}

func (p *Parameter) addImports(pm *PackageMap) {
	p.Type.addImports(pm)
}

// NewParameter returns Parameter.
func NewParameter(name string, t Type) *Parameter {
	return &Parameter{
		Name: name,
		Type: t,
	}
}

// PrintNameAndType print code.
// x int
func (p *Parameter) PrintNameAndType(myPkgPath string, pm PackageMap) string {
	return p.Name + " " + p.Type.PrintDef(myPkgPath, pm)
}

// Contents is contents of source code.
type Contents interface {
	PrintCode(myPkgPath string, pm PackageMap) string
	addImports(pm *PackageMap)
}

// Interface is interface.
type Interface struct {
	Name    string
	Methods []*Method
}

func (i Interface) addImports(pm *PackageMap) {
	for _, m := range i.Methods {
		m.addImports(pm)
	}
}

// PrintCode print code.
func (i *Interface) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		type Foo interface{
			GetXXX (x int, y int) int
		}
	*/
	str := fmt.Sprintf("type %s interface{", i.Name)
	for _, m := range i.Methods {
		str += "\n"
		str += fmt.Sprintf("func%s%s%s", m.Name, m.Type.printArgs(myPkgPath, pm), m.Type.printResults(myPkgPath, pm))
	}
	str += "\n"
	str += "}"
	return str
}

// Struct is struct.
type Struct struct {
	Name    string
	Methods []*Method
	Members []*Member
}

// NewStruct return Struct
func NewStruct(name string) *Struct {
	return &Struct{
		Name:    name,
		Methods: []*Method{},
		Members: []*Member{},
	}
}

func (s *Struct) addImports(pm *PackageMap) {
	for _, m := range s.Methods {
		m.addImports(pm)
	}
	for _, m := range s.Members {
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
	str := fmt.Sprintf("type %s struct{", s.Name)
	for _, m := range s.Members {
		str += "\n"
		str += m.PrintDef(myPkgPath, pm)
	}
	str += "\n"
	str += "}"

	str += "\n"
	for _, m := range s.Methods {
		str += "\n"
		str += m.PrintCode(myPkgPath, pm)
	}
	return str
}

// AddMember add member to struct
func (s *Struct) AddMember(m *Member) {
	s.Members = append(s.Members, m)
}

// AddMethod add method to struct
func (s *Struct) AddMethod(m *Method) {
	s.Methods = append(s.Methods, m)
}

// Member is member of struct.
type Member struct {
	Parameter
	Tag string
}

// NewMember returns Member
func NewMember(name string, t Type, tag string) *Member {
	p := NewParameter(name, t)
	return &Member{
		Parameter: Parameter{
			Name: p.Name,
			Type: p.Type,
		},
		Tag: tag,
	}
}

// PrintDef print name ,type and tag(if exist) .
func (m *Member) PrintDef(myPkgPath string, pm PackageMap) string {
	if m.Tag != "" {
		return fmt.Sprintf("%s `%s`", m.Parameter.PrintNameAndType(myPkgPath, pm), m.Tag)
	}
	return m.Parameter.PrintNameAndType(myPkgPath, pm)
}

// Method is method.
type Method struct {
	Reciever Parameter
	Func
}

// NewMethod returns Method.
func NewMethod(rcv Parameter, name string, t *TypeFunc, body string) *Method {
	f := NewFunc(name, t, body)
	return &Method{
		Reciever: rcv,
		Func: Func{
			Name: f.Name,
			Type: f.Type,
			Body: f.Body,
		},
	}
}

// PrintCode print code.
func (m *Method) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		func (m Reciever) SomeFunc (x int, y int) (int error){
			return x+i , nil
		}
	*/
	s := "func "
	s += fmt.Sprintf("(%s)", m.Reciever.PrintNameAndType(myPkgPath, pm))
	s += m.Name
	s += m.Type.printArgs(myPkgPath, pm)
	s += m.Type.printResults(myPkgPath, pm)
	s += "{\n"

	s += m.Body
	s += "\n"

	s += "}\n"
	return s
}

// Func is func.
type Func struct {
	Name string
	Type *TypeFunc
	Body string
}

// NewFunc returns Func.
func NewFunc(name string, t *TypeFunc, body string) *Func {
	return &Func{
		Name: name,
		Type: t,
		Body: body,
	}
}

// PrintCode print code.
func (f *Func) PrintCode(myPkgPath string, pm PackageMap) string {
	/*
		func SomeFunc (x int, y int) (int error){
			return x+i , nil
		}
	*/
	s := "func "
	s += f.Name
	s += f.Type.printArgs(myPkgPath, pm)
	s += f.Type.printResults(myPkgPath, pm)
	s += "{\n"

	s += f.Body
	s += "\n"

	s += "}\n"
	return s
}

func (f *Func) addImports(pm *PackageMap) {
	f.Type.addImports(pm)
}

// SetBody set body.
func (f *Func) SetBody(body string) {
	f.Body = body
}
