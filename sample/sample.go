package sample

import (
	"codegen/sample/dummy"
	dummy2 "codegen/sample/dummy2"
	"io"

	"github.com/ant0ine/go-json-rest/rest"
)

// SomeInterface is sample.
type SomeInterface interface {
	Foo(arr [2]string, slice []int, n int, chs <-chan string, chr chan<- int, chrs chan int64, ip *int) (map[string]string, error)
	Baa(d dummy.Dummy, d2 dummy2.Dummy2, ss SomeStruct, i io.Writer, si SomeInterface, r *rest.Request)
	Baz(fn func(a int, b int, c dummy.Dummy) error, n int) (a string, b int)
	Qux(
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
}

// NewSomeImplFunc generate SomeImpl
var NewSomeImplFunc = func() (SomeInterface, error) {
	return &SomeImpl{}, nil
}

// SomeImpl is Impliment of SomeInterface
type SomeImpl struct {
	S string
}

// Foo is Foo
func (*SomeImpl) Foo(arr [2]string, slice []int, n int, chs <-chan string, chr chan<- int, chrs chan int64, ip *int) (map[string]string, error) {
	return nil, nil
}

// Baa is Baa
func (*SomeImpl) Baa(d dummy.Dummy, d2 dummy2.Dummy2, ss SomeStruct, i io.Writer, si SomeInterface, r *rest.Request) {

}

// Baz is Baz
func (*SomeImpl) Baz(fn func(a int, b int, c dummy.Dummy) error, n int) (a string, b int) {
	return
}

// Qux is Qux
func (*SomeImpl) Qux(
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
) {

}
