package sample

import (
	"testing"
)

// Test comparable constraint generic interface mock generation
func TestComparerGenericMock(t *testing.T) {
	stub := StubComparer[string]{
		Compare: StubCompare[string]{
			R0: 1,
		},
		Equal: StubEqual[string]{
			R0: true,
		},
	}

	mock := stub.NewMock()

	cmpResult := mock.Compare("a", "b")
	if cmpResult != 1 {
		t.Errorf("Compare() = %v; want %v", cmpResult, 1)
	}

	eqResult := mock.Equal("a", "a")
	if !eqResult {
		t.Errorf("Equal() = %v; want %v", eqResult, true)
	}
}
