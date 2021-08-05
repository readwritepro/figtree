//=============================================================================
// File:     writers_test.go
// Tests:    figtreeWriter with and without includes
//           internalWriter with and without includes
//           jsonWriter with and without includes
//           yamlWriter with and without includes
//=============================================================================

package figtree

import (
	"testing"

	"github.com/joehonton/compare-files"
)

func TestFigtreeWriter(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, _ := ReadConfig(inFilename)

	pw := FigtreeWriter{}
	outFilename := "testdata/actual/sample-figtree"
	err := root.WriteToFile(pw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/sample-figtree", "testdata/actual/sample-figtree")
}

func TestFigtreeIncludes(t *testing.T) {
	inFilename := "testdata/fixtures/include-base"
	root, err := ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	pw := FigtreeWriter{}
	outFilename := "testdata/actual/include-figtree"
	err = root.WriteToFile(pw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/include-figtree", "testdata/actual/include-figtree")
}

func TestInternalWriter(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, _ := ReadConfig(inFilename)

	iw := InternalWriter{}
	outFilename := "testdata/actual/sample-internal"
	err := root.WriteToFile(iw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/sample-internal", "testdata/actual/sample-internal")
}

func TestInternalIncludes(t *testing.T) {
	inFilename := "testdata/fixtures/include-base"
	root, err := ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	iw := InternalWriter{}
	outFilename := "testdata/actual/include-internal"
	err = root.WriteToFile(iw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/include-internal", "testdata/actual/include-internal")
}

func TestJsonWriter(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, _ := ReadConfig(inFilename)

	jw := JsonWriter{}
	outFilename := "testdata/actual/sample-json"
	err := root.WriteToFile(jw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/sample-json", "testdata/actual/sample-json")
}

func TestJsonIncludes(t *testing.T) {
	inFilename := "testdata/fixtures/include-base"
	root, err := ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	jw := JsonWriter{}
	outFilename := "testdata/actual/include-json"
	err = root.WriteToFile(jw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/include-json", "testdata/actual/include-json")
}

func TestYamlWriter(t *testing.T) {
	inFilename := "testdata/fixtures/sample"
	root, _ := ReadConfig(inFilename)

	yw := YamlWriter{}
	outFilename := "testdata/actual/sample-yaml"
	err := root.WriteToFile(yw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/sample-yaml", "testdata/actual/sample-yaml")
}

func TestYamlIncludes(t *testing.T) {
	inFilename := "testdata/fixtures/include-base"
	root, err := ReadConfig(inFilename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	jw := YamlWriter{}
	outFilename := "testdata/actual/include-yaml"
	err = root.WriteToFile(jw, outFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	compare.ExpectedActual(t, "testdata/expected/include-yaml", "testdata/actual/include-yaml")
}
