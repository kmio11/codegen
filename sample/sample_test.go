package sample_test

import (
	"codegen/sample"
	"testing"
)

func TestSample(t *testing.T) {
	stub := sample.StubSomeInterface{
		Name: sample.StubName{
			R0: "mocked",
		},
		Add: sample.StubAdd{
			R0: 100,
			R1: nil,
		},
	}
	stub.SetMock()

	// test
	mock := sample.NewSomeImpl("foo")

	add, err := mock.Add(1, 1)
	if 100 != add || err != nil {
		t.Fatalf("not mocked. Add returns %v", add)
	}

	name := mock.Name()
	if "mocked" != name {
		t.Fatalf("not mocked. Name returns %s", name)
	}
}
