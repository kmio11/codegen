package parser

import (
	"codegen/pkg/generator/model"
	"fmt"
	"go/types"
)

func (p *Parser) parseInterface(obj types.Object) (*model.Interface, error) {
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
		typeFunc, ok := mtype.(*model.TypeFunc)
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
					Type: typeFunc,
				},
				Reciever: model.Parameter{
					Name: "",
					Type: rcv,
				},
			})
	}

	intf := &model.Interface{
		Name:    obj.Name(),
		Methods: methods,
	}
	return intf, nil
}

func (*Parser) getMethodSet(obj types.Object, pointer bool) (*types.MethodSet, error) {
	t := obj.Type()
	if pointer {
		t = types.NewPointer(t)
	}
	return types.NewMethodSet(t), nil
}
