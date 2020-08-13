package model

import "fmt"

// File contains the information related to the file.
type File struct {
	Path         string
	PkgName      string
	PkgPath      string
	Dependencies *PackageMap // packages this package imports.
	Contents     []Contents  // functions, interfaces, structs , ...
}

// NewFile returns File.
func NewFile(path, pkgname, pkgpath string, dependencies *PackageMap) *File {
	return &File{
		Path:         path,
		PkgName:      pkgname,
		PkgPath:      pkgpath,
		Dependencies: dependencies,
	}
}

// Print returns code.
func (f *File) Print() string {
	var s string
	s += fmt.Sprintf("package %s", f.PkgName)
	s += "\n"
	s += f.Dependencies.PrintCode(f.PkgPath)

	for _, c := range f.Contents {
		s += "\n"
		s += c.PrintCode(f.PkgPath, *f.Dependencies)
	}

	return s
}

// ImportsTidy set Dependencies in this file.
func (f *File) ImportsTidy() *PackageMap {
	f.Dependencies.CleanDependencies()
	for _, c := range f.Contents {
		c.addImports(f.Dependencies)
	}

	f.Dependencies.ResolveNameConflict(f.PkgPath)
	return f.Dependencies
}

// AddInterface add interface to file.
func (f *File) AddInterface(intf *Interface) {
	f.Contents = append(f.Contents, intf)
}

// AddStruct add interface to file.
func (f *File) AddStruct(s *Struct) {
	f.Contents = append(f.Contents, s)
}
