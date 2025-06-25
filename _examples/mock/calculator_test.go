package mock

import (
	"errors"
	"testing"
)

func TestCalculator_Mock(t *testing.T) {
	mock := &MockCalculator{
		FakeAdd: func(a, b int) int {
			return a + b
		},
		FakeDivide: func(a, b int) (int, error) {
			if b == 0 {
				return 0, errors.New("division by zero")
			}
			return a / b, nil
		},
	}

	// Test Add
	result := mock.Add(5, 3)
	if result != 8 {
		t.Errorf("Expected 8, got %d", result)
	}

	// Test Divide
	result, err := mock.Divide(10, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}

	// Test Divide by zero
	_, err = mock.Divide(10, 0)
	if err == nil {
		t.Error("Expected error for division by zero")
	}
}

func TestCalculator_Stub(t *testing.T) {
	stub := StubCalculator{
		Add:      StubAdd{R0: 15},
		Subtract: StubSubtract{R0: 5},
		Multiply: StubMultiply{R0: 50},
		Divide:   StubDivide{R0: 2, R1: nil},
	}

	calc := stub.NewMock()

	if result := calc.Add(10, 5); result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}

	if result := calc.Subtract(10, 5); result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}

	if result := calc.Multiply(10, 5); result != 50 {
		t.Errorf("Expected 50, got %d", result)
	}

	result, err := calc.Divide(10, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestStorage_Generic(t *testing.T) {
	// Test string-int storage
	stub := StubStorage[string, int]{
		Get: StubGet[string, int]{
			R0: 42,
			R1: true,
		},
	}

	storage := stub.NewMock()
	storage.Set("key", 42)

	value, found := storage.Get("key")
	if !found {
		t.Error("Expected key to be found")
	}
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}