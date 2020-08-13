package test

import (
	"codegen/test/dummy"
	dummy2 "codegen/test/dummy2"
	"io"

	"github.com/ant0ine/go-json-rest/rest"
)

// SomeInterface is test.
type SomeInterface interface {
	FuncA(arr [2]string, slice []int, n int, chs <-chan string, chr chan<- int, chrs chan int64, ip *int, v ...string) (map[string]string, error)
	FuncB(d dummy.Dummy, d2 dummy2.Dummy2, ss SomeStruct, i io.Writer, si SomeInterface, r *rest.Request)
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
