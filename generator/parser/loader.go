package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// A Package contains all the information related to a parsed package.
type Package struct {
	Name  string
	Files []*ast.File
	Pkg   *types.Package
}

// LoadPackage parse package.
func (p *Parser) LoadPackage(patterns ...string) error {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedModule,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		p.log.Println(err)
		return err
	}
	err = p.PrintErrors(pkgs)
	if err != nil {
		p.log.Println(err)
		return err
	}

	if len(pkgs) != 1 {
		err = fmt.Errorf("error: %d packages found", len(pkgs))
		p.log.Println(err.Error())
		return err
	}

	pkg := pkgs[0]
	p.ParsedPkg = &Package{
		Name:  pkg.Name,
		Pkg:   pkg.Types,
		Files: pkg.Syntax,
	}
	return nil
}

// PrintErrors prints to logger the accumulated errors of all
// packages in the import graph rooted at pkgs, dependencies first.
// PrintErrors returns error if  1 or more errors printed and stopLoadErr is true.
func (p *Parser) PrintErrors(pkgs []*packages.Package) error {
	errs := p.PkgErrors(pkgs)
	if len(errs) > 0 {
		for _, err := range errs {
			p.log.Printf("[WARN] %s\n", err)
		}
		if p.stopLoadErr {
			err := fmt.Errorf("error occuerd when loading package")
			return err
		}
		p.log.Println("[WARN] error occuerd when loading package")
	}
	return nil
}

// PkgErrors returns the accumulated errors of all
// packages in the import graph rooted at pkgs, dependencies first.
func (p *Parser) PkgErrors(pkgs []*packages.Package) []error {
	errs := []error{}
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			errs = append(errs, err)
		}
	})
	return errs
}
