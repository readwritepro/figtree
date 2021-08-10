//=============================================================================
// File:     reader_test.go
// Tests:    Read success
//           Read missing input file
//           Read premature closing brace
//           Read unmatched opening brace
//=============================================================================

package figtree_test

import (
	"fmt"
	"testing"

	"github.com/readwritepro/figtree"
)

func TestRead(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, err := figtree.ReadConfig(inFilename)

	var expectedErr error // nil
	if err != expectedErr {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}

	expected := "*figtree.Branch"
	typeOfRoot := fmt.Sprintf("%T", root)
	if typeOfRoot != expected {
		t.Errorf("expected type '%s', got '%T'", expected, root)
	}
}

func TestMissingInputFile(t *testing.T) {
	inFilename := "testdata/fixtures/missing-config"
	_, err := figtree.ReadConfig(inFilename)

	expectedMsg := fmt.Sprintf("open %s: no such file or directory", inFilename)
	if expectedMsg != err.Error() {
		t.Errorf("expected '%s', got '%v'", expectedMsg, err)
	}
}

func TestPrematureClosingBrace(t *testing.T) {
	inFilename := "testdata/fixtures/premature-closing-brace"
	_, err := figtree.ReadConfig(inFilename)

	expectedErr := figtree.ErrEndOfBranch
	if expectedErr != err {
		t.Errorf("expected '%v', got '%v'", expectedErr, err)
	}
}

func TestUnmatchedOpeningBrace(t *testing.T) {
	inFilename := "testdata/fixtures/unmatched-opening-brace"
	_, err := figtree.ReadConfig(inFilename)

	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}
}
