package parser

import (
	"fmt"
	"go/types"

	"github.com/kmio11/codegen/generator/model"
)

func (p *Parser) parseType(t types.Type) (model.Type, error) {
	var m model.Type
	var err error

	switch tt := t.(type) {
	case *types.Array:
		m, err = p.parseArray(tt)
	case *types.Slice:
		m, err = p.parseSlice(tt)
	case *types.Basic:
		m, err = p.parseBasic(tt)
	case *types.Chan:
		m, err = p.parseChan(tt)
	case *types.Interface:
		m, err = p.parseInterfaceTmp(tt)
	case *types.Map:
		m, err = p.parseMap(tt)
	case *types.Named:
		m, err = p.parseNamedType(tt)
	case *types.Pointer:
		m, err = p.parsePointer(tt)
	case *types.Signature:
		m, err = p.parseSignature(tt)
	case *types.Struct:
		m, err = p.parseStruct(tt)
	default:
		err = fmt.Errorf("unexpected type: %s", tt.String())
	}

	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *Parser) parseArray(t *types.Array) (model.Type, error) {
	tt, err := p.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &model.TypeArray{
		Len:  t.Len(),
		Type: tt,
	}, nil
}

func (p *Parser) parseSlice(t *types.Slice) (model.Type, error) {
	tt, err := p.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &model.TypeArray{
		Len:  -1,
		Type: tt,
	}, nil
}

func (*Parser) parseBasic(t *types.Basic) (model.Type, error) {
	b := model.TypeBasic(t.Name())
	return &b, nil
}

func (p *Parser) parseChan(t *types.Chan) (model.Type, error) {
	tt, err := p.parseType(t.Elem())
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

	return &model.TypeChan{
		Dir:  dir,
		Type: tt,
	}, nil
}

func (p *Parser) parseInterfaceTmp(t *types.Interface) (model.Type, error) {
	embeddeds := []model.Type{}
	for i := 0; i < t.NumEmbeddeds(); i++ {
		embedded, err := p.parseType(t.EmbeddedType(i))
		if err != nil {
			return nil, err
		}
		embeddeds = append(embeddeds, embedded)
	}

	emethods := []*model.Func{}
	for i := 0; i < t.NumExplicitMethods(); i++ {
		tf := t.ExplicitMethod(i)
		emethod, err := p.parseFunc(tf)
		if err != nil {
			return nil, err
		}

		emethods = append(emethods, emethod)
	}

	return &model.TypeInterface{
		Embeddeds:       embeddeds,
		ExplicitMethods: emethods,
	}, nil
}

func (p *Parser) parseMap(t *types.Map) (model.Type, error) {
	k, err := p.parseType(t.Key())
	if err != nil {
		return nil, err
	}
	v, err := p.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &model.TypeMap{
		Key:   k,
		Value: v,
	}, nil
}

func (*Parser) parseNamedType(t *types.Named) (model.Type, error) {
	var pkg *model.PkgInfo
	if t.Obj().Pkg() != nil {
		pkg = model.NewPkgInfo(t.Obj().Pkg().Name(), t.Obj().Pkg().Path(), "")
	}

	return model.NewTypeNamed(pkg, t.Obj().Name()), nil
}

func (p *Parser) parsePointer(t *types.Pointer) (model.Type, error) {
	tt, err := p.parseType(t.Elem())
	if err != nil {
		return nil, err
	}
	return model.NewPointer(tt), nil
}

func (p *Parser) parseSignature(t *types.Signature) (model.Type, error) {
	params, err := p.getParameter(t.Params())
	if err != nil {
		return nil, err
	}
	var variadic *model.Parameter
	if t.Variadic() {
		// If variadic is set, the function is variadic,
		// *types.Signature must have at least one parameter,
		// and the last parameter must be of unnamed slice type.
		lastParam := params[len(params)-1]
		slice, ok := lastParam.Type.(*model.TypeArray)
		if !ok {
			err = fmt.Errorf("internal error. %s is %T", lastParam.Name, lastParam.Type)
			p.log.Println(err)
			return nil, err
		}
		variadic = &model.Parameter{
			Name: lastParam.Name,
			Type: slice.Type,
		}
		params = params[:len(params)-1]
	}

	results, err := p.getParameter(t.Results())
	if err != nil {
		return nil, err
	}

	return &model.TypeSignature{
		Params:   params,
		Variadic: variadic,
		Results:  results,
	}, nil
}

func (p *Parser) parseStruct(t *types.Struct) (model.Type, error) {
	fields := []*model.Field{}
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		name := f.Name()
		if f.Embedded() {
			name = ""
		}
		typ, err := p.parseType(f.Type())
		if err != nil {
			return nil, err
		}
		fields = append(fields, model.NewField(name, typ, t.Tag(i)))
	}

	return &model.TypeStruct{
		Fields: fields,
	}, nil
}

func (p *Parser) getParameter(t *types.Tuple) ([]*model.Parameter, error) {
	params := []*model.Parameter{}
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		vv, err := p.parseType(v.Type())
		if err != nil {
			return nil, err
		}
		params = append(params, &model.Parameter{
			Name: v.Name(),
			Type: vv,
		})
	}
	return params, nil
}

func (p *Parser) parseTuple(t *types.Tuple) ([]model.Type, error) {
	ms := []model.Type{}
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		vv, err := p.parseType(v.Type())
		if err != nil {
			return nil, err
		}
		ms = append(ms, vv)
	}
	return ms, nil
}
