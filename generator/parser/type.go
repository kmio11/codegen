package parser

import (
	"fmt"
	"go/types"

	"github.com/kmio11/codegen/generator/model"
)

func (p *Parser) parseType(t types.Type) (model.Type, error) {
	tp := p.newTypeParser()
	return tp.parseType(t)
}

type typeParser struct {
	p      *Parser
	parsed map[parsedKey]model.Type  // for cycle detection
	stats  map[parsedKey]parseStatus // for cycle detection
}

func (p *Parser) newTypeParser() *typeParser {
	return &typeParser{
		p:      p,
		parsed: map[parsedKey]model.Type{},
		stats:  map[parsedKey]parseStatus{},
	}
}

type parsedKey string

func (tp *typeParser) getParsedKey(t types.Type) parsedKey {
	var key string
	switch tt := t.(type) {
	case *types.Named:
		var path string
		if tt.Obj().Pkg() != nil {
			path = tt.Obj().Pkg().Path()
		}
		key = path + "::" + tt.Obj().Name()
	}

	return parsedKey(key)
}

type parseStatus int

const (
	pstatUnknown parseStatus = iota
	pstatMarked              // you should parse
	pstatParsing             // the othrer is now parsing
	pstatParsed              // already parsed
)

// wasParsed returns parsing status.
// If status is parsing, returns pointer that will be set value when parsing process finish.
func (tp *typeParser) wasParsed(t types.Type) (addr model.Type, stat parseStatus, add func(model.Type) model.Type) {
	key := tp.getParsedKey(t)
	stat = pstatUnknown

	add = func(m model.Type) model.Type {
		ret := m
		if stat == pstatMarked {
			switch mm := m.(type) {
			case *model.TypeNamed:
				p, ok := tp.parsed[key].(*model.TypeNamed)
				if !ok {
					panic("type must be *model.TypeNamed")
				}
				*p = *mm
				ret = p
			default:
				msg := fmt.Sprintf("unexpected type. %T", m)
				panic(msg)
			}
			tp.stats[key] = pstatParsed
		}
		return ret
	}

	switch t.(type) {
	case *types.Named:
		stat = tp.stats[key]
		if stat == pstatUnknown {
			stat = pstatMarked
			tp.parsed[key] = model.NewTypeNamed(nil, "memory allocation", nil)
		}
	}

	if stat == pstatMarked {
		tp.stats[key] = pstatParsing
	}

	if stat >= pstatParsing {
		addr = tp.parsed[key]
	}

	return addr, stat, add
}

func (tp *typeParser) parseType(t types.Type) (model.Type, error) {
	var m model.Type
	var err error

	m, stat, end := tp.wasParsed(t)
	if stat >= pstatParsing {
		return m, nil
	}

	switch tt := t.(type) {
	case *types.Array:
		m, err = tp.parseArray(tt)
	case *types.Slice:
		m, err = tp.parseSlice(tt)
	case *types.Basic:
		m, err = tp.parseBasic(tt)
	case *types.Chan:
		m, err = tp.parseChan(tt)
	case *types.Interface:
		m, err = tp.parseInterface(tt)
	case *types.Map:
		m, err = tp.parseMap(tt)
	case *types.Named:
		m, err = tp.parseNamedType(tt)
	case *types.Pointer:
		m, err = tp.parsePointer(tt)
	case *types.Signature:
		m, err = tp.parseSignature(tt)
	case *types.Struct:
		m, err = tp.parseStruct(tt)
	case *types.TypeParam:
		m, err = tp.parseTypeParam(tt)
	case *types.Union:
		m, err = tp.parseUnionConstraint(tt)
	default:
		err = fmt.Errorf("unexpected type: %s", tt.String())
	}

	if err != nil {
		return nil, err
	}
	m = end(m)
	return m, nil
}

func (tp *typeParser) parseArray(t *types.Array) (model.Type, error) {
	tt, err := tp.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return model.NewTypeArray(t.Len(), tt), nil
}

func (tp *typeParser) parseSlice(t *types.Slice) (model.Type, error) {
	tt, err := tp.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return model.NewTypeArray(-1, tt), nil
}

func (tp *typeParser) parseBasic(t *types.Basic) (model.Type, error) {
	return model.NewTypeBasic(t.Name()), nil
}

func (tp *typeParser) parseChan(t *types.Chan) (model.Type, error) {
	tt, err := tp.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	tdir := t.Dir()
	var dir model.ChanDir
	switch tdir {
	case types.SendRecv:
		dir = model.SendRecv
	case types.SendOnly:
		dir = model.SendOnly
	case types.RecvOnly:
		dir = model.RecvOnly
	}

	return model.NewTypeChan(dir, tt), nil
}

func (tp *typeParser) parseInterface(t *types.Interface) (model.Type, error) {
	embeddeds := []*model.TypeNamed{}
	for i := 0; i < t.NumEmbeddeds(); i++ {
		embedded, err := tp.parseType(t.EmbeddedType(i))
		if err != nil {
			return nil, err
		}
		embeddeds = append(embeddeds, embedded.(*model.TypeNamed))
	}

	emethods := []*model.Func{}
	for i := 0; i < t.NumExplicitMethods(); i++ {
		tf := t.ExplicitMethod(i)
		emethod, err := tp.parseFunc(tf)
		if err != nil {
			return nil, err
		}

		emethods = append(emethods, emethod)
	}

	// Note: For now, we create regular interfaces here
	// Type parameters are handled at the Named type level
	return model.NewTypeInterface(embeddeds, emethods), nil
}

func (tp *typeParser) parseMap(t *types.Map) (model.Type, error) {
	k, err := tp.parseType(t.Key())
	if err != nil {
		return nil, err
	}
	v, err := tp.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return model.NewTypeMap(k, v), nil
}

func (tp *typeParser) parseNamedType(t *types.Named) (model.Type, error) {
	var pkg *model.PkgInfo
	if t.Obj().Pkg() != nil {
		pkg = model.NewPkgInfo(t.Obj().Pkg().Name(), t.Obj().Pkg().Path(), "")
	}
	org, err := tp.parseType(t.Obj().Type().Underlying())
	if err != nil {
		return nil, err
	}

	// Parse type parameters if this is a generic type (Go 1.18+)
	var typeParams []*model.TypeParameter
	if t.TypeParams() != nil && t.TypeParams().Len() > 0 {
		typeParams, err = tp.parseTypeParameters(t.TypeParams())
		if err != nil {
			return nil, err
		}
		return model.NewGenericTypeNamed(pkg, t.Obj().Name(), org, typeParams), nil
	}

	return model.NewTypeNamed(pkg, t.Obj().Name(), org), nil
}

func (tp *typeParser) parsePointer(t *types.Pointer) (model.Type, error) {
	tt, err := tp.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return model.NewPointer(tt), nil
}

func (tp *typeParser) parseSignature(t *types.Signature) (model.Type, error) {
	params, err := tp.parseParameter(t.Params())
	if err != nil {
		return nil, err
	}
	var variadic *model.Parameter
	if t.Variadic() {
		// If variadic is set, the function is variadic,
		// *types.Signature must have at least one parameter,
		// and the last parameter must be of unnamed slice type.
		lastParam := params[len(params)-1]
		slice, ok := lastParam.Type().(*model.TypeArray)
		if !ok {
			err = fmt.Errorf("internal error. %s is %T", lastParam.Name(), lastParam.Type())
			tp.p.log.Println(err)
			return nil, err
		}
		variadic = model.NewParameter(lastParam.Name(), slice.Type())
		params = params[:len(params)-1]
	}

	results, err := tp.parseParameter(t.Results())
	if err != nil {
		return nil, err
	}

	return model.NewTypeSignature(params, variadic, results), nil
}

func (tp *typeParser) parseStruct(t *types.Struct) (model.Type, error) {
	fields := []*model.Field{}
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		name := f.Name()
		if f.Embedded() {
			name = ""
		}
		typ, err := tp.parseType(f.Type())
		if err != nil {
			return nil, err
		}
		fields = append(fields, model.NewField(name, typ, t.Tag(i)))
	}

	return model.NewTypeStruct(fields), nil
}

func (tp *typeParser) parseParameter(t *types.Tuple) ([]*model.Parameter, error) {
	params := []*model.Parameter{}
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		typ, err := tp.parseType(v.Type())
		if err != nil {
			return nil, err
		}
		params = append(params, model.NewParameter(v.Name(), typ))
	}
	return params, nil
}

func (tp *typeParser) parseTuple(t *types.Tuple) ([]model.Type, error) {
	ms := []model.Type{}
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		vv, err := tp.parseType(v.Type())
		if err != nil {
			return nil, err
		}
		ms = append(ms, vv)
	}
	return ms, nil
}

func (tp *typeParser) parseFunc(t *types.Func) (*model.Func, error) {
	sig, _ := t.Type().(*types.Signature) // *types.Func's Type() is always a *Signature
	typ, err := tp.parseType(sig)
	if err != nil {
		return nil, err
	}
	typf, ok := typ.(*model.TypeSignature)
	if !ok {
		err = fmt.Errorf("internal error. not TypeSignature: <%s>'s type is <%T> ", t.String(), typf)
		tp.p.log.Println(err)
		return nil, err
	}

	return model.NewFunc(t.Name(), typf, ""), nil
}

// parseTypeParameters parses type parameter list
func (tp *typeParser) parseTypeParameters(typeParams *types.TypeParamList) ([]*model.TypeParameter, error) {
	var params []*model.TypeParameter
	for i := 0; i < typeParams.Len(); i++ {
		param := typeParams.At(i)
		constraint, err := tp.parseTypeConstraint(param.Constraint())
		if err != nil {
			return nil, err
		}
		typeParam := model.NewTypeParameter(param.Obj().Name(), constraint, i)
		params = append(params, typeParam)
	}
	return params, nil
}

// parseTypeConstraint parses type constraint
func (tp *typeParser) parseTypeConstraint(constraint types.Type) (model.Type, error) {
	switch c := constraint.(type) {
	case *types.Interface:
		// Handle built-in constraints and interface constraints
		if c.IsComparable() {
			return model.ConstraintComparable, nil
		}
		if c.Empty() {
			return model.ConstraintAny, nil
		}
		// For complex interface constraints, parse as regular interface
		return tp.parseInterface(c)
	case *types.Union:
		// Handle union constraints like `int | string`
		return tp.parseUnionConstraint(c)
	case *types.Basic:
		// Handle basic type constraints
		return tp.parseBasic(c)
	case *types.Named:
		// Handle named type constraints
		return tp.parseNamedType(c)
	default:
		// Default to any for unknown constraints
		return model.ConstraintAny, nil
	}
}

// parseUnionConstraint parses union type constraints
func (tp *typeParser) parseUnionConstraint(union *types.Union) (model.Type, error) {
	// For now, we'll represent union constraints as a string
	// This could be enhanced later to support proper union types
	var constraintName string
	for i := 0; i < union.Len(); i++ {
		term := union.Term(i)
		if i > 0 {
			constraintName += " | "
		}
		constraintName += term.Type().String()
	}
	return model.NewTypeConstraint(constraintName), nil
}

// parseTypeParam parses a single type parameter
func (tp *typeParser) parseTypeParam(param *types.TypeParam) (model.Type, error) {
	// For type parameters that appear in method signatures, we represent them as their name
	constraint, err := tp.parseTypeConstraint(param.Constraint())
	if err != nil {
		return nil, err
	}
	return model.NewTypeParameter(param.Obj().Name(), constraint, param.Index()), nil
}
