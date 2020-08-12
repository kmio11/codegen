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
	Name  string
	Path  string
	Alias string
}

// NewPkgInfo returns ImportedPackage
func NewPkgInfo(name, path, alias string) *PkgInfo {
	return &PkgInfo{
		Name:  name,
		Path:  path,
		Alias: alias,
	}
}

// PrintCode returns import declaration.
func (i *PkgInfo) PrintCode() string {
	if i.Alias != "" {
		return fmt.Sprintf(`%s "%s"`, i.Alias, i.Path)
	}
	return fmt.Sprintf(`"%s"`, i.Path)
}

// Prefix returns prefix to use types in myPkg.
func (i *PkgInfo) Prefix(myPkg string) string {
	if myPkg == i.Path {
		return ""
	}

	if i.Alias == dotImport {
		return ""
	}
	if i.Alias != "" {
		return i.Alias + "."
	}
	return i.Name + "."
}

// PackageMap is packages.
type PackageMap struct {
	list     map[string]PkgInfo
	isImport map[string]bool
}

// NewPackageMap returns PackageMap
func NewPackageMap(myPkgName, myPkgPath string) *PackageMap {
	p := &PackageMap{
		list:     map[string]PkgInfo{},
		isImport: map[string]bool{},
	}
	p.Add(myPkgPath, *NewPkgInfo(myPkgName, myPkgPath, ""))
	return p
}

// Copy copy PackageMap to dst.
func (pm *PackageMap) Copy(dst *PackageMap) {
	dst.list = map[string]PkgInfo{}
	dst.isImport = map[string]bool{}

	for k, v := range pm.list {
		dst.list[k] = v
	}
	for k, v := range pm.isImport {
		dst.isImport[k] = v
	}
}

// MarshalJSON is marshal json
func (pm *PackageMap) MarshalJSON() ([]byte, error) {
	// sort by pkgpath
	paths := []string{}
	for k := range pm.list {
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
			Name:       pm.list[path].Name,
			Alias:      pm.list[path].Alias,
			IsImported: pm.isImport[path],
		})
	}
	return json.Marshal(outs)
}

// CleanDependencies set all packages as  isImported=false.
func (pm *PackageMap) CleanDependencies() {
	for k := range pm.isImport {
		pm.isImport[k] = false
	}
}

// needToImport returns package paths need to import.
func (pm *PackageMap) needToImport(myPkgPath string) []string {
	paths := []string{}
	for path, need := range pm.isImport {
		if myPkgPath == path || !need {
			continue
		}
		paths = append(paths, path)
	}

	sort.Strings(paths)
	return paths
}

// PrintCode print code.
func (pm *PackageMap) PrintCode(myPkgPath string) string {
	str := ""

	paths := pm.needToImport(myPkgPath)
	if len(paths) == 0 {
		return str
	}

	str = "import ("
	for _, path := range paths {
		pkg := pm.list[path]
		str += "\n"
		str += pkg.PrintCode()
	}
	str += "\n"
	str += ")"
	return str
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

	paths := pm.needToImport(myPkgPath)
	used := make([]string, len(paths))
	for n, path := range paths {
		var imName string
		pkginfo := pm.list[path]
		if pkginfo.Alias == dotImport || pkginfo.Alias == brankImport {
			continue
		}
		if pkginfo.Alias != "" {
			imName = pkginfo.Alias
		} else {
			imName = pkginfo.Name
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
			pkginfo.Alias = imName
			pm.list[path] = pkginfo
		}
	}
}

// Need set package as isImported=need
func (pm *PackageMap) Need(path string, need bool) {
	pm.isImport[path] = need
}

// Get returns ImportedPackage
func (pm *PackageMap) Get(path string) *PkgInfo {
	// return pm.list[path]
	v, ok := pm.list[path]
	if !ok {
		return nil
	}
	return &v
}

// Add add Package to map.
func (pm *PackageMap) Add(path string, pkg PkgInfo) {
	pm.list[path] = pkg
	pm.isImport[path] = true
}
