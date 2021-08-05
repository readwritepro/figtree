//=============================================================================
// File:   alter_test.go
// Tests:  Branch.Append
//         Branch.InsertBeforeItem
//         Branch.InsertAfterItem
//         Branch.RemoveItem
//         FigtreeWriter.WriteToBuffer
//=============================================================================

package figtree

import (
	"testing"
)

func TestAppendItem(t *testing.T) {
	root := NewBranch()

	keyval1 := Item{
		key:   "key1",
		value: "value1",
	}

	root.AppendItem(keyval1)
	actual, _ := root.GetValue("/key1")

	expected := "value1"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}

func TestInsertBeforeItem(t *testing.T) {
	root := NewBranch()

	keyval0 := Item{
		key:   "key0",
		value: "value0",
	}
	keyval1 := Item{
		key:   "key1",
		value: "value1",
	}

	keyval2 := Item{
		key:   "key2",
		value: "value2",
	}

	// exercise the "not found" case
	err := root.InsertBeforeItem("key2", keyval0)
	if err != ErrNotFound {
		t.Errorf("expected '%v', got '%v'", ErrNotFound, err)
	}

	// add items 2, 0, 1
	root.AppendItem(keyval2)
	root.InsertBeforeItem("key2", keyval0) // insert before only item
	root.InsertBeforeItem("key2", keyval1) // insert before last item, after first item

	// result should be ordered 0, 1, 2
	fw := FigtreeWriter{}
	actual, _ := root.WriteToBuffer(fw)
	expected := "key0 value0\nkey1 value1\nkey2 value2\n"

	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}

func TestInsertAfterItem(t *testing.T) {
	root := NewBranch()

	keyval0 := Item{
		key:   "key0",
		value: "value0",
	}
	keyval1 := Item{
		key:   "key1",
		value: "value1",
	}

	keyval2 := Item{
		key:   "key2",
		value: "value2",
	}

	// exercise the "not found" case
	err := root.InsertAfterItem("key0", keyval1)
	if err != ErrNotFound {
		t.Errorf("expected '%v', got '%v'", ErrNotFound, err)
	}

	// add items 0, 1, 2
	root.AppendItem(keyval0)
	root.InsertAfterItem("key0", keyval2) // insert after only item
	root.InsertAfterItem("key0", keyval1) // insert after first item, before second item

	// result should be ordered 0, 1, 2
	fw := FigtreeWriter{}
	actual, _ := root.WriteToBuffer(fw)
	expected := "key0 value0\nkey1 value1\nkey2 value2\n"

	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}

func TestRemoveItem(t *testing.T) {
	root := NewBranch()

	keyval0 := Item{
		key:   "key0",
		value: "value0",
	}
	keyval1 := Item{
		key:   "key1",
		value: "value1",
	}

	keyval2 := Item{
		key:   "key2",
		value: "value2",
	}

	// exercise the "not found" case
	err := root.RemoveItem("key0")
	if err != ErrNotFound {
		t.Errorf("expected '%v', got '%v'", ErrNotFound, err)
	}

	// add items 0, 1, 2
	root.AppendItem(keyval0)
	root.AppendItem(keyval1)
	root.AppendItem(keyval2)

	root.RemoveItem("key0") // remove first item, leaving 1, 2
	root.RemoveItem("key2") // remove last item, leaving 1
	root.RemoveItem("key1") // remove last item, leaving nothing

	// result should be ordered nothing
	fw := FigtreeWriter{}
	actual, _ := root.WriteToBuffer(fw)
	expected := ""

	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}
