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
	methods := []*model.Func{}
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
		methods = append(methods,
			model.NewFunc(
				method.Obj().Name(),
				sig,
				"",
			),
		)
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
