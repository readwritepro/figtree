//=============================================================================
// File:     merge_test.go
// Tests:    Read user file with !baseline pragma
//           Sort items
//           Write merged baseline + user file
//=============================================================================

package figtree

import (
	"testing"

	"github.com/joehonton/compare-files"
)

func TestUserWithBaselineFile(t *testing.T) {
	inFilename := "testdata/fixtures/user"
	root, err := ReadConfig(inFilename)
	if err != nil {
		t.Fatal()
	}

	root.SortItems()

	iw := InternalWriter{}
	outFilename := "testdata/actual/user-internal"
	err = root.WriteToFile(iw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	compare.ExpectedActual(t, "testdata/expected/user-internal", "testdata/actual/user-internal")

	fw := FigtreeWriter{}
	outFilename = "testdata/actual/user-figtree"
	err = root.WriteToFile(fw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	compare.ExpectedActual(t, "testdata/expected/user-figtree", "testdata/actual/user-figtree")
}
