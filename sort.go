//=============================================================================
// File:     sort.go
// Contents: Sort the items of a branch
//=============================================================================

package figtree

import "sort"

// returns a string value that can be used to sort this item amongst its siblings
func (item Item) sortKey() string {
	switch item.value.(type) {
	case string:
		return "0" + item.key + " " + item.value.(string)
	case *Branch:
		return "1" + item.key
	}
	return ""
}

// Reorder the items in the branch with key/value items first, followed by inner branches.
// Items at each level are sorted alphabetically by keyName.
// This is a recursive function.
func (branch *Branch) SortItems() {
	sort.Slice(branch.Items, func(i, j int) bool {
		return branch.Items[i].sortKey() < branch.Items[j].sortKey()
	})

	for _, item := range branch.Items {
		switch value := item.value.(type) {
		case *Branch:
			value.SortItems()
		default:
		}
	}
}
