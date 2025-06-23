package model

import (
	"fmt"
	"strings"
)

// PrintType returns type string.
func printType(typ Type, myPkgPath string, pm PackageMap) string {
	switch t := typ.(type) {
	case *TypeParameter:
		return t.name

	case *TypeConstraint:
		return t.name

	case *TypeArray:
		s := "[]"
		if t.Len() > -1 {
			s = fmt.Sprintf("[%d]", t.Len())
		}
		return s + t.Type().PrintType(myPkgPath, pm)

	case *TypeBasic:
		return string(*t)

	case *TypeChan:
		s := t.Type().PrintType(myPkgPath, pm)
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

	case *TypeInterface:
		s := "interface{"

		// Add type parameters if this is a generic interface
		if len(t.typeParams) > 0 {
			s = "interface["
			for i, param := range t.typeParams {
				if i > 0 {
					s += ", "
				}
				s += param.name
				if param.constraint != nil {
					s += " " + param.constraint.PrintType(myPkgPath, pm)
				}
			}
			s += "]{"
		}

		for _, e := range t.Embeddeds() {
			s += e.PrintType(myPkgPath, pm)
			s += ";"
		}
		for _, e := range t.ExplicitMethods() {
			s += e.PrintDef(myPkgPath, pm)
			s += ";"
		}
		s = strings.TrimRight(s, ";")
		s += "}"

		return s

	case *TypeMap:
		return fmt.Sprintf("map[%s]%s", t.Key().PrintType(myPkgPath, pm), t.Value().PrintType(myPkgPath, pm))

	case *TypeNamed:
		var result string
		if t.Pkg() == nil {
			result = t.Name()
		} else {
			pkg := pm.Get(t.Pkg().Path())
			if pkg == nil {
				pkg = t.Pkg()
			}
			result = pkg.Prefix(myPkgPath) + t.Name()
		}

		// Add type parameters if this is a generic type
		if len(t.typeParams) > 0 {
			result += "["
			for i, param := range t.typeParams {
				if i > 0 {
					result += ", "
				}
				result += param.name
				// Don't print constraints when used as type arguments in fields
				// Only print constraints in type definitions
			}
			result += "]"
		}

		return result

	case *TypePointer:
		return "*" + t.Type().PrintType(myPkgPath, pm)

	case *TypeSignature:
		s := "func"
		s += t.printArgs(myPkgPath, pm)
		s += t.printResults(myPkgPath, pm)
		return s

	case *TypeStruct:
		s := "struct {\n"
		for _, f := range t.Fields() {
			s += "\t" + f.PrintDef(myPkgPath, pm) + "\n"
		}
		s += "}"

		return s

	default:
		panic("unexpected type")
	}
}

// printArgs print params
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

// printResults print params
// for example : (x int, y int)
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
