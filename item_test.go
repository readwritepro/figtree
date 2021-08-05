//=============================================================================
// File:     item_test.go
// Tests:    Type, Key, SetKey, Value, SetValue, Branch, SetBranch
//=============================================================================

package figtree

import "testing"

func TestItem(t *testing.T) {
	var actual, actualType, expected string
	var actualErr, expectedErr error

	item := Item{
		key:   "key1",
		value: "value1",
	}

	// change the key
	item.SetKey("new key")
	actual = item.Key()
	expected = "new key"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	// change the value
	item.SetValue("new value")
	actual, _ = item.Value()
	expected = "new value"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	// change type from [leaf] to [branch]
	item.SetBranch(NewBranch())
	actualType = item.Type()
	expected = "[branch]"
	if expected != actualType {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	// attempt to incorrectly use the item as a leaf
	actual, actualErr = item.Value()
	expectedErr = ErrNotLeaf
	if expectedErr != actualErr {
		t.Errorf("expected '%v', got '%v'", expectedErr, actualErr)
	}

	// change type from [branch] to [leaf]
	item.SetValue("value2")
	actualType = item.Type()
	expected = "[leaf]"
	if expected != actualType {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	// attempt to incorrectly use the item as a branch
	_, actualErr = item.Branch()
	expectedErr = ErrNotBranch
	if expectedErr != ErrNotBranch {
		t.Errorf("expected '%v', got '%v'", expectedErr, actualErr)
	}

}
