package model

// typeImports add packages requrired by typ to pm.
func typeImports(typ Type, pm *PackageMap) {
	switch t := typ.(type) {
	case *TypeArray:
		t.Type().addImports(pm)

	case *TypeBasic:
		// do nothing

	case *TypeChan:
		t.Type().addImports(pm)

	case *TypeInterface:
		for _, e := range t.Embeddeds() {
			e.addImports(pm)
		}
		for _, e := range t.ExplicitMethods() {
			e.addImports(pm)
		}

	case *TypeMap:
		t.Key().addImports(pm)
		t.Value().addImports(pm)

	case *TypeNamed:
		if t.Pkg() != nil {
			pm.SetRequired(t.pkg.Path(), true)
		}

	case *TypePointer:
		t.Type().addImports(pm)

	case *TypeSignature:
		for _, p := range t.Args() {
			p.addImports(pm)
		}
		if t.Variadic() != nil {
			t.Variadic().addImports(pm)
		}
		for _, p := range t.Results() {
			p.addImports(pm)
		}

	case *TypeStruct:
		for _, f := range t.Fields() {
			f.addImports(pm)
		}

	default:
		panic("unexpected type")
	}
}
