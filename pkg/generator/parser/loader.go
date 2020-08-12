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
	e := packages.PrintErrors(pkgs)
	if e > 0 && p.stopLoadErr {
		err := fmt.Errorf("error occuerd when loading package")
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
