package parser

import (
	"codegen/pkg/generator/model"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
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
		log: log.New(ioutil.Discard, "", 0),
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

	} else {
		return fmt.Errorf("%s is unsupported", obj.Type())
	}
	return nil
}
