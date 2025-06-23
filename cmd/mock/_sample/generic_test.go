package sample

import (
	"testing"
)

// Test type that implements Stringable interface
type TestStringable struct {
	value string
}

func (ts TestStringable) String() string {
	return ts.value
}

// Test basic generic interface mock generation
func TestRepositoryGenericMock(t *testing.T) {
	// This test will verify that we can generate mocks for Repository[T any]
	// Expected: MockRepository[T any] should be generated with proper type parameters

	// Test with string type
	stubStr := StubRepository[string]{
		Get: StubGet[string]{
			R0: "test-item",
			R1: nil,
		},
		Save: StubSave[string]{
			R0: nil,
		},
		List: StubList[string]{
			R0: []string{"item1", "item2"},
			R1: nil,
		},
	}

	mockStr := stubStr.NewMock()

	// Test Get method
	item, err := mockStr.Get("test-id")
	if item != "test-item" || err != nil {
		t.Errorf("Get() = %v, %v; want %v, %v", item, err, "test-item", nil)
	}

	// Test Save method
	err = mockStr.Save("new-item")
	if err != nil {
		t.Errorf("Save() error = %v; want nil", err)
	}

	// Test List method
	items, err := mockStr.List()
	if len(items) != 2 || items[0] != "item1" || err != nil {
		t.Errorf("List() = %v, %v; want %v, %v", items, err, []string{"item1", "item2"}, nil)
	}

	// Test with int type
	stubInt := StubRepository[int]{
		Get: StubGet[int]{
			R0: 42,
			R1: nil,
		},
	}

	mockInt := stubInt.NewMock()
	intItem, err := mockInt.Get("test-id")
	if intItem != 42 || err != nil {
		t.Errorf("Get() = %v, %v; want %v, %v", intItem, err, 42, nil)
	}
}

// Test multiple type parameters
func TestConverterGenericMock(t *testing.T) {
	stub := StubConverter[string, int]{
		Convert: StubConvert[string, int]{
			R0: 123,
			R1: nil,
		},
	}

	mock := stub.NewMock()
	result, err := mock.Convert("test")
	if result != 123 || err != nil {
		t.Errorf("Convert() = %v, %v; want %v, %v", result, err, 123, nil)
	}
}

// Test comparable constraint
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

// Test union type constraint
func TestNumberGenericMock(t *testing.T) {
	// Test with int
	stubInt := StubNumber[int]{
		Add: StubAdd[int]{
			R0: 5,
		},
		Multiply: StubMultiply[int]{
			R0: 20,
		},
	}

	mockInt := stubInt.NewMock()

	addResult := mockInt.Add(2, 3)
	if addResult != 5 {
		t.Errorf("Add() = %v; want %v", addResult, 5)
	}

	mulResult := mockInt.Multiply(4, 5)
	if mulResult != 20 {
		t.Errorf("Multiply() = %v; want %v", mulResult, 20)
	}

	// Test with float64
	stubFloat := StubNumber[float64]{
		Add: StubAdd[float64]{
			R0: 5.5,
		},
	}

	mockFloat := stubFloat.NewMock()
	floatResult := mockFloat.Add(2.2, 3.3)
	if floatResult != 5.5 {
		t.Errorf("Add() = %v; want %v", floatResult, 5.5)
	}
}

// Test interface constraint
func TestFormatterGenericMock(t *testing.T) {
	stub := StubFormatter[TestStringable]{
		Format: StubFormat[TestStringable]{
			R0: "formatted",
		},
		FormatList: StubFormatList[TestStringable]{
			R0: []string{"fmt1", "fmt2"},
		},
	}

	mock := stub.NewMock()

	item := TestStringable{value: "test"}
	formatted := mock.Format(item)
	if formatted != "formatted" {
		t.Errorf("Format() = %v; want %v", formatted, "formatted")
	}

	items := []TestStringable{{value: "a"}, {value: "b"}}
	formattedList := mock.FormatList(items)
	if len(formattedList) != 2 || formattedList[0] != "fmt1" {
		t.Errorf("FormatList() = %v; want %v", formattedList, []string{"fmt1", "fmt2"})
	}
}
