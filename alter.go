//=============================================================================
// File:     alter.go
// Contents: Functions to alter the figtree hierarchy
//           AppendItem, appendItem, PrependItem
//           InsertBeforeItem, InsertAfterItem
//           RemoveItem
//=============================================================================

package figtree

// Add an item to the branch, appending it to the end of the branch's list of key/values.
func (branch *Branch) AppendItem(item Item) {
	branch.Items = append(branch.Items, item)
}

func (branch *Branch) appendItem(key string, value interface{}, blockComments []string, terminalWhitespace string, terminalComment string, srcFile string, srcLine *int, srcOrigin FileOrigin) {
	item := Item{
		key:                key,
		value:              value,
		blockComments:      blockComments,
		terminalWhitespace: terminalWhitespace,
		terminalComment:    terminalComment,
		srcFile:            srcFile,
		srcLine:            *srcLine,
		srcOrigin:          srcOrigin,
	}
	branch.Items = append(branch.Items, item)
}

// Add an item to the branch, prepending it to the beginning of the branch's list of key/values.
func (branch *Branch) PrependItem(item Item) {
	newBranchItems := make([]Item, len(branch.Items)+1) // make room for the new item with a larger slice
	newBranchItems[0] = item
	for j := 0; j < len(branch.Items); j++ {
		newBranchItems[j+1] = (branch.Items)[j]
	}
	branch.Items = newBranchItems
}

// Add an item to the current branch, placing it immediately before the targetKeyName.
//
// Returns ErrNotFound if the targetKeyName is not in the current branch.
func (branch *Branch) InsertBeforeItem(targetKeyName string, newItem Item) error {
	for index, item := range branch.Items {
		if item.key == targetKeyName {
			newBranchItems := make([]Item, len(branch.Items)+1) // make room for the new item with a larger slice
			for j := 0; j < index; j++ {
				newBranchItems[j] = (branch.Items)[j]
			}
			newBranchItems[index] = newItem
			for j := index; j < len(branch.Items); j++ {
				newBranchItems[j+1] = (branch.Items)[j]
			}
			branch.Items = newBranchItems
			return nil
		}
	}
	return ErrNotFound
}

// Add an item to the current branch, placing it immediately after the targetKeyName.
//
// Returns ErrNotFound if the targetKeyName is not in the current branch.
func (branch *Branch) InsertAfterItem(targetKeyName string, newItem Item) error {
	for index, item := range branch.Items {
		if item.key == targetKeyName {
			newBranchItems := make([]Item, len(branch.Items)+1) // make room for the new item with a larger slice
			for j := 0; j < index+1; j++ {
				newBranchItems[j] = (branch.Items)[j]
			}
			newBranchItems[index+1] = newItem
			for j := index + 1; j < len(branch.Items); j++ {
				newBranchItems[j+1] = (branch.Items)[j]
			}
			branch.Items = newBranchItems
			return nil
		}
	}
	return ErrNotFound
}

// Remove the item with the specified keyName from the given branch.
// When a branch has more than one item with such a keyName, only the first one is removed.
//
// Returns ErrNotFound if the specified keyName is not in the current branch.
func (branch *Branch) RemoveItem(keyName string) error {
	for index, item := range branch.Items {
		if item.key == keyName {
			branch.Items = append(branch.Items[:index], branch.Items[index+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
