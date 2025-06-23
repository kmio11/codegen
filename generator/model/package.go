package model

// A Package contains all the information related to a parsed package.
type Package struct {
	Name         string
	Path         string
	Dependencies *PackageMap // packages this package imports
	Interfaces   []*Interface
	Functions    []*Func
	Structs      []*Struct
}

// CopyDependencies return copy of Dependencies.
func (p *Package) CopyDependencies() *PackageMap {
	pm := NewPackageMap(p.Name, p.Path)
	p.Dependencies.Copy(pm)
	return pm
}
