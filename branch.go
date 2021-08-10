//=============================================================================
// File:     branch.go
// Contents: Branch type declaration
//           CreateBranch
//=============================================================================

package figtree

// The Branch type is a slice of configuration tree items, in insertion order, where items
// are either key/value pairs or key/branch pairs. Branches form a hierarchical tree of items.
type Branch struct {
	Items []Item
}

// The NewBranch function is used to create a branch that will be used
// as a starting point for in-memory figtree manipulation.
// It may also be used when creating an inner branch that will be added with one
// of the alter functions (AppendItem, PrependItem, InsertBeforeItem, InsertAfterItem).
func NewBranch() *Branch {
	branch := Branch{}
	return &branch
}

var gBaselineTree *Branch // pointer to the tree of items built by the !baseline pragma, if any
var gDtdTree *Branch      // pointer to the tree of items built by the !dtd pragma, if any
