//=============================================================================
// File:     main_test.go
// Contents: SETUP and TEARDOWN for tests
//=============================================================================

package figtree_test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("=== SETUP")
	os.Remove("testdata/actual/include-figtree")
	os.Remove("testdata/actual/include-internal")
	os.Remove("testdata/actual/sample-figtree")
	os.Remove("testdata/actual/sample-internal")

	rc := m.Run()

	fmt.Println("=== TEARDOWN")
	os.Exit(rc)
}
