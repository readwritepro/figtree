//=============================================================================
// File:     branch.go
// Contents: Branch type declaration
//           CreateBranch
//=============================================================================

// Package figtree provides a multi-paradigm SDK for sophisticated configuration
// file access.
//
// Figtree syntax is based on classic key/value pairs that may be organized
// into a nested hierarchy of named sections.
// Many of the design goals for the figtree syntax come from its predecessors
// including XML, JSON, YAML, TOML, win.ini, and Apache.
//
// There are two essential constructs to understand. First, simple key/value
// pairs appear on a single line. Unlike other configuration syntax styles,
// key/value pairs do not use colons or equal-signs as assignment operators.
// Instead, the presence of whitespace between the key-name and the beginning
// of the value is enough to parse the line into left-hand and right-hand halves.
// Furthermore, values are not delimited by quotation marks. Instead, all leading
// and trailing whitespace characters are stripped from the right-hand side of
// the key/value pair to produce the value.
//
// The second construct to understand are named sections, which are multi-line
// collections of key/value pairs. Named sections have a key and a pair of braces.
// The opening brace must be the last non-whitespace character of a line. The
// closing brace must be the first non-whitespace character of a line. The lines
// between the opening and closing braces may contain simple key/value pairs
// or other named sections, which may be nested arbitrarily deep.
//
// Block comments are written using hashtags as the first non-whitespace character
// of a line. Terminal comments may appear together with key/value pairs when
// whitespace and a hash tag follow the value.
//
package figtree

// The Branch type is a slice of configuration tree items, in insertion order.
type Branch struct {
	Items      []Item
	branchPath string // The full keyPathName from the root branch down to this branch
}

// The NewBranch function is used to create a branch that will be used
// as a starting point for in-memory figtree manipulation.
// It may also be used when creating an inner branch that will be added with one
// of the alter functions (AppendItem, InsertBeforeItem, InsertAfterItem).
func NewBranch() *Branch {
	branch := Branch{}
	return &branch
}

var gBaselineTree *Branch // pointer to the tree of items built by the !baseline pragma, if any
var gDtdTree *Branch      // pointer to the tree of items built by the !dtd pragma, if any
