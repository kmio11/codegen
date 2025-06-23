package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

var goReservedWords = []string{
	"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan",
	"else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue",
	"for", "import", "return", "var",
}

const (
	dotImport   = "."
	brankImport = "_"
)

// PkgInfo is package infomation.
type PkgInfo struct {
	name  string
	path  string
	alias string
}

// NewPkgInfo returns ImportedPackage
func NewPkgInfo(name, path, alias string) *PkgInfo {
	return &PkgInfo{
		name:  name,
		path:  path,
		alias: alias,
	}
}

// Name returns package name.
func (pi *PkgInfo) Name() string {
	return pi.name
}

// Path returns package path.
func (pi *PkgInfo) Path() string {
	return pi.path
}

// Alias returns package alias.
func (pi *PkgInfo) Alias() string {
	return pi.alias
}

// PrintCode returns import declaration.
func (pi *PkgInfo) PrintCode() string {
	if pi.alias != "" {
		return fmt.Sprintf(`%s "%s"`, pi.alias, pi.path)
	}
	return fmt.Sprintf(`"%s"`, pi.path)
}

// Prefix returns prefix to use types in myPkg.
func (pi *PkgInfo) Prefix(myPkg string) string {
	if myPkg == pi.path {
		return ""
	}

	if pi.alias == dotImport {
		return ""
	}
	if pi.alias != "" {
		return pi.alias + "."
	}
	return pi.name + "."
}

// PackageMap is packages.
type PackageMap struct {
	pkgs    map[string]PkgInfo
	imports map[string]bool
}

// NewPackageMap returns PackageMap
func NewPackageMap(myPkgName, myPkgPath string) *PackageMap {
	p := &PackageMap{
		pkgs:    map[string]PkgInfo{},
		imports: map[string]bool{},
	}
	p.Add(myPkgPath, *NewPkgInfo(myPkgName, myPkgPath, ""))
	return p
}

// PrintCode print code.
func (pm *PackageMap) PrintCode(myPkgPath string) string {
	str := ""

	paths := pm.requireImport(myPkgPath)
	if len(paths) == 0 {
		return str
	}

	str = "import ("
	for _, path := range paths {
		pkg := pm.pkgs[path]
		str += "\n"
		str += pkg.PrintCode()
	}
	str += "\n"
	str += ")"
	return str
}

// MarshalJSON is marshal json
func (pm *PackageMap) MarshalJSON() ([]byte, error) {
	// sort by pkgpath
	paths := []string{}
	for k := range pm.pkgs {
		paths = append(paths, k)
	}

	type out struct {
		Path       string
		Name       string
		Alias      string
		IsImported bool
	}
	outs := []out{}
	for _, path := range paths {
		outs = append(outs, out{
			Path:       path,
			Name:       pm.pkgs[path].name,
			Alias:      pm.pkgs[path].alias,
			IsImported: pm.imports[path],
		})
	}
	return json.Marshal(outs)
}

// Copy copy PackageMap to dst.
func (pm *PackageMap) Copy(dst *PackageMap) {
	dst.pkgs = map[string]PkgInfo{}
	dst.imports = map[string]bool{}

	for k, v := range pm.pkgs {
		dst.pkgs[k] = v
	}
	for k, v := range pm.imports {
		dst.imports[k] = v
	}
}

// CleanDependencies set all packages to not require import.
func (pm *PackageMap) CleanDependencies() {
	for k := range pm.imports {
		pm.imports[k] = false
	}
}

// requireImport returns package paths myPkgPath need to import.
func (pm *PackageMap) requireImport(myPkgPath string) []string {
	paths := []string{}
	for path, require := range pm.imports {
		if myPkgPath == path || !require {
			continue
		}
		paths = append(paths, path)
	}

	sort.Strings(paths)
	return paths
}

// ResolveNameConflict set alias to packages which need to imported  if name is duplicated.
func (pm *PackageMap) ResolveNameConflict(myPkgPath string) {
	contains := func(s []string, str string) bool {
		for _, v := range s {
			if str == v {
				return true
			}
		}
		return false
	}

	paths := pm.requireImport(myPkgPath)
	used := make([]string, len(paths))
	for n, path := range paths {
		var imName string
		pkginfo := pm.pkgs[path]
		if pkginfo.alias == dotImport || pkginfo.alias == brankImport {
			continue
		}
		if pkginfo.alias != "" {
			imName = pkginfo.alias
		} else {
			imName = pkginfo.name
		}

		var isAliasUpdate bool
		tmp := imName
		for i := 0; ; i++ {
			if contains(used, imName) || contains(goReservedWords, imName) {
				isAliasUpdate = true
				imName = tmp + strconv.Itoa(i)
				continue
			}
			break
		}
		used[n] = imName
		if isAliasUpdate {
			pkginfo.alias = imName
			pm.pkgs[path] = pkginfo
		}
	}
}

// SetRequired set package as isImported=need
func (pm *PackageMap) SetRequired(path string, isRequred bool) {
	pm.imports[path] = isRequred
}

// Get returns ImportedPackage
func (pm *PackageMap) Get(path string) *PkgInfo {
	// return pm.list[path]
	v, ok := pm.pkgs[path]
	if !ok {
		return nil
	}
	return &v
}

// Add add package, and set it to requred to import.
func (pm *PackageMap) Add(path string, pkg PkgInfo) {
	pm.pkgs[path] = pkg
	pm.imports[path] = true
}
