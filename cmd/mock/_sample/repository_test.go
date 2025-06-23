package sample

import (
	"testing"
)

// Test basic generic interface mock generation for Repository[T any]
func TestRepositoryGenericMock(t *testing.T) {
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
