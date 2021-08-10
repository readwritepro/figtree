package figtree_test

import (
	"fmt"
	"os"

	"github.com/readwritepro/figtree"
)

func ExampleReadConfig() {
	inFilename := "testdata/fixtures/example-config"
	root, err := figtree.ReadConfig(inFilename)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Unable to read config file %s", inFilename)
	}
	fmt.Fprintf(os.Stdout, "%T", root)

	// Output:
	// *figtree.Branch
}
