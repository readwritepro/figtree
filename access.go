//=============================================================================
// File:     access.go
// Contents: Functions for querying and retreiving the figtree hierarchy
//           QueryAll, QueryOne
//           ListBranches, ListLeaves, ItemCount
//           GetItem, GetBranch, GetLeaf, GetValue
//           ItemExists, ItemIsBranch, ItemIsLeaf, ItemIsArray, PathExists
//=============================================================================

package figtree

import (
	"strings"
)

// Find all items with the given key path.
// This method provides access to multiple items functioning like an array.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// Returns a collection of zero or more items that fully match the keyPath.
// The items in the collection may contain a mixture of both leaves (strings) and Branches.
func (branch *Branch) QueryAll(keyPath string) []Item {
	// is this a simpleKeyName or a pathKeyName
	slash := strings.Index(keyPath, "/")

	// recurse on "/"
	if slash == 0 {
		rightSide := keyPath[1:]
		return branch.QueryAll(rightSide)
	}

	// recurse on pathKeyName
	if slash != -1 {
		// split the path at the first branch/key boundary
		leftSide := keyPath[:slash]
		rightSide := keyPath[slash+1:]

		// lookup the inner branch
		for _, item := range branch.Items {
			if item.key == leftSide {
				innerBranch := item.value.(*Branch)
				return innerBranch.QueryAll(rightSide)
			}
		}
		return make([]Item, 0)
	}

	// return all leaves that fully match simpleKeyName
	collection := make([]Item, 0)
	for _, item := range branch.Items {
		if item.key == keyPath {
			collection = append(collection, item)
		}
	}
	return collection
}

// Find the first item with the given key path.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// Returns a single item that fully match the keyPath.
// If more than one match exists, returns the first one found.
// If no match exists returns a nil Item pointer and ErrNotFound.
func (branch *Branch) QueryOne(keyPath string) (*Item, error) {
	// is this a simpleKeyName or a pathKeyName
	slash := strings.Index(keyPath, "/")

	// recurse on "/"
	if slash == 0 {
		rightSide := keyPath[1:]
		collection := branch.QueryAll(rightSide)
		if len(collection) == 0 {
			return nil, ErrNotFound
		} else {
			return &collection[0], nil
		}
	}

	// recurse on pathKeyName
	if slash != -1 {
		// split the path at the first branch/key boundary
		leftSide := keyPath[:slash]
		rightSide := keyPath[slash+1:]

		// lookup the inner branch
		for _, item := range branch.Items {
			if item.key == leftSide {
				innerBranch := item.value.(*Branch)
				collection := innerBranch.QueryAll(rightSide)
				if len(collection) == 0 {
					return nil, ErrNotFound
				} else {
					return &collection[0], nil
				}
			}
		}
		return nil, ErrNotFound
	}

	// return the first leaf that fully matches simpleKeyName
	for _, item := range branch.Items {
		if item.key == keyPath {
			return &item, nil
		}
	}
	return nil, ErrNotFound
}

// Get all sub-branches of the current branch
// Returns a collection of zero or more subordinate branches
func (branch *Branch) ListBranches() []Item {
	var subBranches []Item

	for _, item := range branch.Items {
		if item.Type() == "[branch]" {
			subBranches = append(subBranches, item)
		}
	}

	return subBranches
}

// Get all leaves of the current branch
// Returns a collection of zero or more leaf items
func (branch *Branch) ListLeaves() []Item {
	var leaves []Item

	for _, item := range branch.Items {
		if item.Type() == "[leaf]" {
			leaves = append(leaves, item)
		}
	}

	return leaves
}

// Get the item with the given keyPath.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// Returns ErrNotFound when the keyPath does not exist.
// Returns ErrNotLeaf if the keyPath points to a leaf rather than a branch.
func (branch *Branch) GetItem(keyPath string) (*Item, error) {
	// is this a simpleKeyName or a pathKeyName
	slash := strings.Index(keyPath, "/")

	// recurse on "/"
	if slash == 0 {
		rightSide := keyPath[1:]
		pItem, err := branch.GetItem(rightSide)
		if err != nil {
			return nil, err
		} else {
			return pItem, nil
		}
	}

	// recurse on pathKeyName
	if slash != -1 {
		// split the path at the first branch/key boundary
		leftSide := keyPath[:slash]
		rightSide := keyPath[slash+1:]

		pItem, err := branch.GetItem(leftSide)
		if err != nil {
			return nil, err
		}

		innerBranch, err := pItem.Branch()
		if err != nil {
			return nil, err
		}
		return innerBranch.GetItem(rightSide)
	}

	// lookup the current branch
	for index := range branch.Items {
		if branch.Items[index].key == keyPath {
			return &branch.Items[index], nil
		}
	}
	return nil, ErrNotFound
}

// Get the branch with the given keyPath.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// Returns ErrNotFound when the keyPath does not exist.
// Returns ErrNotBranch if the keyPath points to a leaf rather than a branch.
func (branch *Branch) GetBranch(keyPath string) (*Branch, error) {
	item, err := branch.QueryOne(keyPath)
	if err != nil {
		return nil, err
	}

	switch value := item.value.(type) {
	case string:
		return nil, ErrNotBranch // last part of the Path is a leaf
	case *Branch:
		return value, nil // last part of the Path is a *Branch
	default:
		return nil, ErrNotBranch
	}
}

// Get the leaf with the given keyPath.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// Returns ErrNotFound when the keyPath does not exist.
// Returns ErrNotLeaf if the keyPath points to a branch rather than a leaf.
func (branch *Branch) GetLeaf(keyPath string) (*Item, error) {
	item, err := branch.QueryOne(keyPath)
	if err != nil {
		return nil, err
	}

	switch item.value.(type) {
	case string:
		return item, nil // last part of the Path is a leaf
	case *Branch:
		return nil, ErrNotLeaf // last part of the Path is a *Branch
	default:
		return nil, ErrNotLeaf
	}
}

// Find the value of the item with the given keyPath.
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
//
// The returned string is the value corresponding to the key.
// It will be an empty string when it is a "key-only" item.
// Returns ErrNotFound when the keyPath does not exist.
// Returns ErrNotLeaf when the keyPath is a branch rather than a leaf.
func (branch *Branch) GetValue(keyPath string) (string, error) {
	item, err := branch.QueryOne(keyPath)
	if err != nil {
		return "", err
	}

	switch value := item.value.(type) {
	case string:
		return value, nil // last part of the Path is a leaf
	case *Branch:
		return "", ErrNotLeaf // last part of the Path is a *Branch
	default:
		return "", err
	}
}

// Determines whether an item with the given simpleKeyName exists in this branch.
// This function only accepts simpleKeyNames, not keyPaths.
func (branch *Branch) ItemExists(simpleKeyName string) bool {
	for _, item := range branch.Items {
		if item.key == simpleKeyName {
			return true
		}
	}
	return false
}

// Determines whether an item is a branch
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
func (branch *Branch) ItemIsBranch(keyPath string) bool {
	item, err := branch.QueryOne(keyPath)
	if err != nil {
		return false
	}

	switch item.value.(type) {
	case string:
		return false
	case *Branch:
		return true
	default:
		return false
	}
}

// Determines whether an item is a leaf
// The keyPath argument may be a "simpleKeyName" or a "keyPath"
// (a keyPath containing a slash-separated prefix of branch names).
func (branch *Branch) ItemIsLeaf(keyPath string) bool {
	item, err := branch.QueryOne(keyPath)
	if err != nil {
		return false
	}

	switch item.value.(type) {
	case string:
		return true
	case *Branch:
		return false
	default:
		return false
	}
}

// Returns true if this branch has more than one item with the given simpleKeyName.
// This function only accepts simpleKeyNames, not keyPaths.
func (branch *Branch) ItemIsArray(simpleKeyName string) bool {
	count := 0
	for _, item := range branch.Items {
		if item.key == simpleKeyName {
			count++
		}
	}
	return count > 1
}

// Determines whether an item with the given keyPath exists in this branch.
// This is a recursive function that searches down successively deeper branches
// of a slash-separated keyPath.
func (branch *Branch) PathExists(keyPath string) bool {

	// discard leading "/" and recurse
	slash := strings.Index(keyPath, "/")
	if slash == 0 {
		rightSide := keyPath[1:]
		return branch.PathExists(rightSide)
	}

	// split in half and recurse
	if slash != -1 {
		leftSide := keyPath[:slash]
		rightSide := keyPath[slash+1:]
		for _, item := range branch.Items {
			if item.key == leftSide {
				innerBranch := item.value.(*Branch)
				return innerBranch.PathExists(rightSide)
			}
		}
		return false
	}

	// return the first leaf that fully matches
	return branch.ItemExists(keyPath)
}
