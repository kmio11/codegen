package test

import (
	"io"

	"github.com/kmio11/codegen/test/dummy"
	dummy2 "github.com/kmio11/codegen/test/dummy2"
)

// SomeInterface is test.
type SomeInterface interface {
	FuncA(arr [2]string, slice []int, n int, chs <-chan string, chr chan<- int, chrs chan int64, ip *int, v ...string) (map[string]string, error)
	FuncB(d dummy.Dummy, d2 dummy2.Dummy2, ss SomeStruct, i io.Writer, si SomeInterface, r *dummy.Dummy)
	FuncC(fn func(a int, b int, c dummy.Dummy) error, n int) (a string, b int)
	FuncD(
		i interface{},
		ii interface {
			SomeInterface
			Func(x int) int
		},
		s struct{},
		ss struct {
			dummy.Dummy
			em  dummy2.Dummy2
			str string
		},
	)
	FuncE(string, string, int)
}

type CycleIntf interface {
	Cycle(cc CycleStr)
	FuncB(d SomeStruct, r SomeStruct)
}

type CycleStr struct {
	Cs CycleIntf
}
