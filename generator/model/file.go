package model

import "fmt"

// File contains the information related to the file.
type File struct {
	path         string
	pkg          *PkgInfo
	dependencies *PackageMap
	contents     []Contents
}

// NewFile returns File.
func NewFile(path, pkgname, pkgpath string, dependencies *PackageMap) *File {
	return &File{
		path:         path,
		pkg:          NewPkgInfo(pkgname, pkgpath, ""),
		dependencies: dependencies,
	}
}

// Path returns filepath.
func (f *File) Path() string {
	return f.path
}

// Pkg returns PkgInfo the file belongs to.
func (f *File) Pkg() *PkgInfo {
	return f.pkg
}

// Dependencies returns packages the file may depends on.
func (f *File) Dependencies() *PackageMap {
	return f.dependencies
}

// DependenciesTidy add missing and remove unused package.
func (f *File) DependenciesTidy() *PackageMap {
	f.dependencies.CleanDependencies()
	for _, c := range f.contents {
		c.addImports(f.dependencies)
	}

	f.dependencies.ResolveNameConflict(f.pkg.Path())
	return f.dependencies
}

// Contents returnss contents of file. functions, interfaces, structs , ...
func (f *File) Contents() []Contents {
	return f.contents
}

// PrintCode returns code.
func (f *File) PrintCode() string {
	var s string
	s += fmt.Sprintf("package %s", f.pkg.Name())
	s += "\n"
	s += f.dependencies.PrintCode(f.pkg.Path())

	for _, c := range f.contents {
		s += "\n"
		s += c.PrintCode(f.pkg.Path(), *f.dependencies)
	}

	return s
}

// AddInterface add interface to file.
func (f *File) AddInterface(intf *Interface) {
	f.contents = append(f.contents, intf)
}

// AddStruct add interface to file.
func (f *File) AddStruct(s *Struct) {
	f.contents = append(f.contents, s)
}
