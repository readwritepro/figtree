//=============================================================================
// File:     merge_test.go
// Tests:    Read user file with !baseline pragma
//           Sort items
//           Write merged baseline + user file
//=============================================================================

package figtree_test

import (
	"testing"

	"github.com/readwritepro/compare-test-results"
	"github.com/readwritepro/figtree"
)

func TestMerge(t *testing.T) {
	inFilename := "testdata/fixtures/user"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		t.Fatal()
	}

	root.SortItems()

	wi := figtree.WriteInternal{}
	outFilename := "testdata/actual/user-internal"
	err = root.WriteToFile(wi, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	compare.ExpectedActual(t, "testdata/expected/user-internal", "testdata/actual/user-internal")

	wf := figtree.WriteFigtree{}
	outFilename = "testdata/actual/user-figtree"
	err = root.WriteToFile(wf, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	compare.ExpectedActual(t, "testdata/expected/user-figtree", "testdata/actual/user-figtree")
}
