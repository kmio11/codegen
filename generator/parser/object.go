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
			model.NewMethod(
				model.NewParameter("", rcv),
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
