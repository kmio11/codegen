package parser

import (
	"fmt"
	"go/types"

	"github.com/kmio11/codegen/generator/model"
)

func (p *Parser) parseInterfaceObj(obj types.Object) (*model.Interface, error) {
	mset, err := p.getMethodSet(obj, false)
	if err != nil {
		return nil, err
	}
	methods := []*model.Method{}
	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)
		mtype, err := p.parseType(method.Type())
		if err != nil {
			return nil, err
		}
		sig, ok := mtype.(*model.TypeSignature)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}
		rcv, err := p.parseType(method.Recv())
		if err != nil {
			return nil, err
		}
		methods = append(methods,
			&model.Method{
				Func: model.Func{
					Name: method.Obj().Name(),
					Type: sig,
				},
				Reciever: model.Parameter{
					Name: "",
					Type: rcv,
				},
			})
	}

	intf := model.NewInterface(
		obj.Name(),
		model.NewPkgInfo(obj.Pkg().Name(), obj.Pkg().Path(), ""),
		methods,
	)
	return intf, nil
}

func (*Parser) getMethodSet(obj types.Object, pointer bool) (*types.MethodSet, error) {
	t := obj.Type()
	if pointer {
		t = types.NewPointer(t)
	}
	return types.NewMethodSet(t), nil
}

func (p *Parser) parseFunc(t *types.Func) (*model.Func, error) {
	sig, _ := t.Type().(*types.Signature) // *types.Func's Type() is always a *Signature
	typ, err := p.parseSignature(sig)
	if err != nil {
		return nil, err
	}
	typf, ok := typ.(*model.TypeSignature)
	if !ok {
		err = fmt.Errorf("internal error. not TypeSignature: <%s>'s type is <%T> ", t.String(), typf)
		p.log.Println(err)
		return nil, err
	}

	return model.NewFunc(t.Name(), typf, ""), nil
}
