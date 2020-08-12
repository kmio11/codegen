package parser

import (
	"codegen/pkg/generator/model"
	"fmt"
	"go/types"
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
	case *types.Map:
		m, err = p.parseMap(tt)
	case *types.Named:
		m, err = p.parseNamedType(tt)
	case *types.Pointer:
		m, err = p.parsePointer(tt)
	case *types.Signature:
		m, err = p.parseSigature(tt)
	default:
		err = fmt.Errorf("unsupported type: %s", tt.String())
	}

	if err != nil {
		p.log.Println(err)
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

func (p *Parser) parseSigature(t *types.Signature) (model.Type, error) {
	params, err := p.getParameter(t.Params())
	if err != nil {
		return nil, err
	}
	var variadic *model.Parameter
	if t.Variadic() {
		variadic = params[len(params)]
		params = params[:len(params)-1]
	}

	results, err := p.getParameter(t.Results())
	if err != nil {
		return nil, err
	}

	return &model.TypeFunc{
		Params:   params,
		Variadic: variadic,
		Results:  results,
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
