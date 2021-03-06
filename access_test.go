//=============================================================================
// File:   access_test.go
// Tests:  QueryAll
//         QueryOne
//         GetItem
//         GetBranch
//         GetLeaf
//         GetValue
//         ItemExists
//         ItemIsBranch
//         ItemIsLeaf
//         ItemIsArray
//         ItemCount
//         PathExists
//=============================================================================

package figtree_test

import (
	"fmt"
	"testing"

	"github.com/readwritepro/figtree"
)

func TestQueryAll(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	collection := root.QueryAll("/key3")
	actual := len(collection)
	expected := 1
	if expected != actual {
		t.Errorf("expected '%d', got '%d'", expected, actual)
	}

	collection = root.QueryAll("section2/four-identical-keys")
	actual = len(collection)
	expected = 4
	if expected != actual {
		t.Errorf("expected '%d', got '%d'", expected, actual)
	}
}

func TestQueryOne(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	item, _ := root.QueryOne("/key3")
	actual := item.Type()
	expected := "[leaf]"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	item, _ = root.QueryOne("section1")
	actual = item.Type()
	expected = "[branch]"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	item, _ = root.QueryOne("section1/key1")
	actual = item.Type()
	expected = "[leaf]"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	item, _ = root.QueryOne("section2/four-identical-keys")
	actual, _ = item.Value()
	expected = "value1"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	_, err = root.QueryOne("section1/key99")
	expectedErr := figtree.ErrNotFound
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}
func TestListBranches(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	section4, _ := root.GetBranch("section4")
	subBranches := section4.ListBranches()
	count := len(subBranches)
	if count != 2 {
		t.Errorf("expected 2 subBranches, got '%d'", count)
		return
	}

	expected := "section4A"
	actual := subBranches[0].Key()
	if expected != actual {
		t.Errorf("expected '%v', got '%s'", expected, actual)
	}
	expected = "section4B"
	actual = subBranches[1].Key()
	if expected != actual {
		t.Errorf("expected '%v', got '%s'", expected, actual)
	}
}

func TestListLeaves(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	section1, _ := root.GetBranch("section1")
	leaves := section1.ListLeaves()
	count := len(leaves)
	if count != 3 {
		t.Errorf("expected 3 leaves, got '%d'", count)
		return
	}

	expected := "key1"
	actual := leaves[0].Key()
	if expected != actual {
		t.Errorf("expected '%v', got '%s'", expected, actual)
	}
	expected = "key2"
	actual = leaves[1].Key()
	if expected != actual {
		t.Errorf("expected '%v', got '%s'", expected, actual)
	}
	expected = "key3"
	actual = leaves[2].Key()
	if expected != actual {
		t.Errorf("expected '%v', got '%s'", expected, actual)
	}
}

func TestGetItem(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, _ := root.GetItem("section1")
	actualType := fmt.Sprintf("%T", actual)
	expectedType := "*figtree.Item"
	if actualType != expectedType {
		t.Errorf("expected '%v', got '%s'", expectedType, actualType)
	}

	actual, err = root.GetItem("section1/key1")
	actualType = fmt.Sprintf("%T", actual)
	expectedType = "*figtree.Item"
	var expectedErr error = nil
	if actualType != expectedType {
		t.Errorf("expected '%v', got '%s'", expectedType, actualType)
	}
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}

	_, err = root.GetItem("section1/key99")
	expectedErr = figtree.ErrNotFound
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}

func TestGetBranch(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, _ := root.GetBranch("section1")
	actualType := fmt.Sprintf("%T", actual)
	expectedType := "*figtree.Branch"
	if actualType != expectedType {
		t.Errorf("expected '%v', got '%s'", expectedType, actualType)
	}

	_, err = root.GetBranch("section1/key1")
	expectedErr := figtree.ErrNotBranch
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}

	_, err = root.GetBranch("section1/key99")
	expectedErr = figtree.ErrNotFound
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}

func TestGetLeaf(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, _ := root.GetLeaf("section1/key1")
	actualType := fmt.Sprintf("%T", actual)
	expectedType := "*figtree.Item"
	if actualType != expectedType {
		t.Errorf("expected '%v', got '%s'", expectedType, actualType)
	}

	_, err = root.GetLeaf("section1")
	expectedErr := figtree.ErrNotLeaf
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}

	_, err = root.GetLeaf("section1/key99")
	expectedErr = figtree.ErrNotFound
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}

func TestGetValue(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, _ := root.GetValue("/key3")
	expected := "value3"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	_, err = root.GetValue("section1")
	expectedErr := figtree.ErrNotLeaf
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}

	actual, _ = root.GetValue("section1/key1")
	expected = "space-then-value"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	actual, _ = root.GetValue("section2/four-identical-keys")
	expected = "value1"
	if expected != actual {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	_, err = root.GetValue("section1/key99")
	expectedErr = figtree.ErrNotFound
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}

func TestItemIsBranch(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual := root.ItemIsBranch("key1")
	expected := false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemIsBranch("section1")
	expected = true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemIsBranch("key99")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestItemIsLeaf(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual := root.ItemIsLeaf("key1")
	expected := true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemIsLeaf("section1")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemIsLeaf("key99")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestItemExists(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual := root.ItemExists("key1")
	expected := true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemExists("section1")
	expected = true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemExists("key99")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestItemIsArray(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual := root.ItemIsArray("key1")
	expected := false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.ItemIsArray("key99")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	section2, _ := root.GetBranch("/section2")
	actual = section2.ItemIsArray("four-identical-keys")
	expected = true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestItemCount(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	section3, _ := root.GetItem("section3")
	expected := 2
	actual := section3.ItemCount()
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	key3, _ := root.GetItem("key3")
	expected = 0
	actual = key3.ItemCount()
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}

func TestPathExists(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual := root.PathExists("key1")
	expected := true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.PathExists("/section1")
	expected = true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.PathExists("/section1/key1")
	expected = true
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}

	actual = root.PathExists("/section1/key99")
	expected = false
	if expected != actual {
		t.Errorf("expected '%v', got '%v'", expected, actual)
	}
}
