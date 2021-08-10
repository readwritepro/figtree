//=============================================================================
// File:     errors.go
// Contents: Error constants and sentinals
//           TODO sink for compiler "declared but not used" messages
//
//=============================================================================

package figtree

// The Error type is declared to allow error instances to be const.
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEOF             = Error("figtree: end of file")
	ErrNotFound        = Error("figtree: item not found")
	ErrNotBranch       = Error("figtree: item is not a branch")
	ErrNotLeaf         = Error("figtree: item is not a leaf")
	ErrUnknownItemType = Error("figtree: unknown Item type")
)

// ErrEndOfBranch is a sentinal returned from the recursive call to parse an inner branch.
// This is a normal return signal for inner branches, but when returned all the
// way to Read, it signals a premature end to parsing, due to an early closing brace,
const ErrEndOfBranch = Error("figtree: end of branch")

// A /dev/null for compiler "declared but not used" messages. This is a development-only
// function to bypass the compilation error caused by a variable never being used.
// When in production, there should be no calls to this function.
func todo(x ...interface{}) {}

func init() {
	todo(gDtdTree)
}
