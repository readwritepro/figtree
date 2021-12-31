//=============================================================================
// File:     item.go
// Contents: Item type declaration
//           Copy constructor
//           Type, Key, SetKey, Value, SetValue, Branch, SetBranch
//=============================================================================

package figtree

// An item holds a key/value pair where the value may be a string or a Branch pointer.
// The struct's key and value are publicly accessible via the Key, SetKey, Value, SetValue,
// Branch, and SetBranch functions.
//
// Private fields
//
// The blockComments field contains any empty lines or block comment lines that immediately preceed this item.
// The terminalWhitespace field is a string containing the tabs and spaces that separate the value and the terminalComment, if any.
// The terminalComment field contains any comment situated on the same line, to the right of the item's value.
// The srcFile, srcLine, and srcOrigin fields reference the source's filename, line number, and type of origin.
type Item struct {
	key                string
	value              interface{}
	blockComments      []string
	terminalWhitespace string
	terminalComment    string
	srcFile            string
	srcLine            int
	srcOrigin          FileOrigin
}

// Allocate and initialize a new item.
func NewItem(key string, value string) Item {
	newItem := Item{
		key:   key,
		value: value,
	}
	return newItem
}

// Make a copy of an item.
func (item Item) Copy() Item {
	newItem := Item{
		key:                item.key,
		value:              item.value,
		blockComments:      item.blockComments,
		terminalWhitespace: item.terminalWhitespace,
		terminalComment:    item.terminalComment,
		srcFile:            item.srcFile,
		srcLine:            item.srcLine,
		srcOrigin:          item.srcOrigin,
	}
	return newItem
}

// Determine whether the item is a leaf or a branch.
//
// Returns "[leaf]" or "[branch]".
func (item Item) Type() string {
	switch item.value.(type) {
	case string:
		return "[leaf]"
	case *Branch:
		return "[branch]"
	}
	return ""
}

// Returns the number of items in the branch
// Returns 0 if the item is a leaf
func (item Item) ItemCount() int {
	switch value := item.value.(type) {
	case string:
		return 0
	case *Branch:
		return value.ItemCount()
	default:
		return 0
	}
}

// Get the item's key.
func (item Item) Key() string {
	return item.key
}

// Change the item's keyName.
func (item *Item) SetKey(keyName string) {
	item.key = keyName
}

// Get the item's value.
//
// Returns ErrNotLeaf if the item holds a branch pointer rather than a leaf value.
func (item Item) Value() (string, error) {
	switch value := item.value.(type) {
	case string:
		return value, nil
	case *Branch:
		return "", ErrNotLeaf
	}
	return "", ErrNotLeaf
}

// Change the item's value.
func (item *Item) SetValue(value string) {
	item.value = value
}

// Get the item's branch pointer.
//
// Returns ErrNotBranch if the item does not point to a branch.
func (item Item) Branch() (*Branch, error) {
	switch value := item.value.(type) {
	case string:
		return nil, ErrNotBranch
	case *Branch:
		return value, nil
	}
	return nil, ErrNotBranch
}

// Changes an item's value to be the specified branch pointer.
func (item *Item) SetBranch(branch *Branch) {
	item.value = branch
}
