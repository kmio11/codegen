package parser

import (
	"fmt"
	"go/types"
	"io"
	"log"

	"github.com/kmio11/codegen/generator/model"
)

const (
	logPrefix = "[parser]"
)

// Parser is parser
// parsing types.Package and generate model.Package.
type Parser struct {
	ParsedPkg   *Package
	Targets     []string // if nil , all element is parsed.
	stopLoadErr bool
	log         *log.Logger
}

// NewParser returns Parser.
func NewParser(opts ...Opts) *Parser {
	p := &Parser{
		log: log.New(io.Discard, "", 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Parse parse and return parsed model.
func (p *Parser) Parse() (*model.Package, error) {
	if p.ParsedPkg == nil {
		err := fmt.Errorf("invalid parser settings")
		p.log.Println(err)
		return nil, err
	}
	if len(p.Targets) != 1 {
		err := fmt.Errorf("unsupported parser settings")
		p.log.Println(err)
		return nil, err
	}

	pkg, err := p.getPackageBase()
	if err != nil {
		p.log.Println(err)
		return nil, err
	}
	for _, tname := range p.Targets {
		err = p.setContents(pkg, tname)
		if err != nil {
			p.log.Println(err)
			return nil, err
		}
	}

	return pkg, nil
}

func (p *Parser) getPackageBase() (*model.Package, error) {
	pp := p.ParsedPkg.Pkg
	pkg := &model.Package{
		Name:         pp.Name(),
		Path:         pp.Path(),
		Dependencies: model.NewPackageMap(pp.Name(), pp.Path()),
	}

	impkgs := pp.Imports()
	for _, impkg := range impkgs {
		pkg.Dependencies.Add(impkg.Path(), *model.NewPkgInfo(impkg.Name(), impkg.Path(), ""))
	}

	return pkg, nil
}

func (p *Parser) setContents(pkg *model.Package, name string) error {
	obj := p.ParsedPkg.Pkg.Scope().Lookup(name)
	if obj == nil {
		return fmt.Errorf("%s not found", name)
	}

	if types.IsInterface(obj.Type()) {
		intf, err := p.parseInterfaceObj(obj)
		if err != nil {
			return err
		}
		pkg.Interfaces = append(pkg.Interfaces, intf)

	} else if isStruct(obj.Type()) {
		intf, err := p.parseStructAsInterface(obj)
		if err != nil {
			return err
		}
		pkg.Interfaces = append(pkg.Interfaces, intf)

	} else {
		return fmt.Errorf("%s is unsupported", obj.Type())
	}
	return nil
}

// isStruct checks if the given type is a struct type
func isStruct(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return true
	case *types.Named:
		_, ok := u.Underlying().(*types.Struct)
		return ok
	}
	return false
}

// parseStructAsInterface converts a struct to an interface by extracting its methods
func (p *Parser) parseStructAsInterface(obj types.Object) (*model.Interface, error) {
	structType := obj.Type()
	
	// Get method set including both value and pointer receiver methods
	methodSet := types.NewMethodSet(structType)
	pointerMethodSet := types.NewMethodSet(types.NewPointer(structType))
	
	// Combine both method sets
	allMethods := make(map[string]*types.Func)
	
	// Add value receiver methods
	for i := 0; i < methodSet.Len(); i++ {
		sel := methodSet.At(i)
		if method, ok := sel.Obj().(*types.Func); ok && method.Exported() {
			allMethods[method.Name()] = method
		}
	}
	
	// Add pointer receiver methods
	for i := 0; i < pointerMethodSet.Len(); i++ {
		sel := pointerMethodSet.At(i)
		if method, ok := sel.Obj().(*types.Func); ok && method.Exported() {
			allMethods[method.Name()] = method
		}
	}
	
	// Convert to model.Func
	var modelMethods []*model.Func
	for _, method := range allMethods {
		sig := method.Type().(*types.Signature)
		
		// Parse parameters
		var params []*model.Parameter
		if sig.Params() != nil {
			for i := 0; i < sig.Params().Len(); i++ {
				param := sig.Params().At(i)
				paramType, err := p.parseType(param.Type())
				if err != nil {
					return nil, fmt.Errorf("failed to parse parameter type for method %s: %v", method.Name(), err)
				}
				params = append(params, model.NewParameter(param.Name(), paramType))
			}
		}
		
		// Parse return values
		var returns []*model.Parameter
		if sig.Results() != nil {
			for i := 0; i < sig.Results().Len(); i++ {
				result := sig.Results().At(i)
				returnType, err := p.parseType(result.Type())
				if err != nil {
					return nil, fmt.Errorf("failed to parse return type for method %s: %v", method.Name(), err)
				}
				returns = append(returns, model.NewParameter(result.Name(), returnType))
			}
		}
		
		// Create type signature
		typeSig := model.NewTypeSignature(params, nil, returns)
		
		// Create model function
		modelMethod := model.NewFunc(method.Name(), typeSig, "")
		modelMethods = append(modelMethods, modelMethod)
	}
	
	// Create interface name by appending "Interface" to struct name
	interfaceName := obj.Name() + "Interface"
	
	// Create package info (path, name, alias)
	pkgInfo := model.NewPkgInfo(obj.Pkg().Name(), obj.Pkg().Path(), "")
	
	// Create interface
	intf := model.NewInterface(interfaceName, pkgInfo, modelMethods)
	
	return intf, nil
}
