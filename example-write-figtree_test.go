//=============================================================================
// File:     example-write-figtree_test.go
//=============================================================================

package figtree_test

import "github.com/readwritepro/figtree"

func ExampleWriteFigtree() {
	inFilename := "testdata/fixtures/example-config"
	root, _ := figtree.ReadConfig(inFilename)

	wf := figtree.WriteFigtree{}
	outFilename := "testdata/actual/example-figtree"
	root.WriteToFile(wf, outFilename)

	// Output:
	//
}
