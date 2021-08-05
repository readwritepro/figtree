//=============================================================================
// File:     merge.go
// Contents: Merge a user config file with baseline file containing fallback defaults
//           Copy constructor for branches
//           Merge function to
//=============================================================================

package figtree

// Merge the baseline tree with the user's tree. The baseline may be nil.
func mergeBaselineWithUser(baselineTree *Branch, userTree *Branch) *Branch {
	if baselineTree == nil {
		return userTree
	} else {
		mergedTree := baselineTree.Copy()
		mergedTree.Merge(userTree)
		return mergedTree
	}
}

// Make a copy of the branch by copying the items of the current branch
func (branch *Branch) Copy() *Branch {
	newBranch := Branch{
		Items:      make([]Item, 0, len(branch.Items)),
		branchPath: branch.branchPath,
	}
	for _, item := range branch.Items {
		newItem := item.Copy()
		newBranch.Items = append(newBranch.Items, newItem)
	}
	return &newBranch
}

// Recursively merge all items in source into destination
// replacing any item that already exists.
// Typically the dstBranch object is the baseline tree containing fallback values
// and the srcBranch is the user's config tree containing explicit overrides.
func (dstBranch *Branch) Merge(srcBranch *Branch) {

	alreadySeen := make(map[string]bool)

	for _, srcItem := range srcBranch.Items {
		key := srcItem.key
		if srcBranch.ItemIsArray(key) || dstBranch.ItemIsArray(key) {
			_, exists := alreadySeen[key]
			if !exists {
				dstBranch.mergeArrayItems(srcBranch, key)
			}
			alreadySeen[key] = true
		} else {
			dstBranch.mergeScalarItem(srcItem)
		}
	}
}

// Merge the srcItem into the dstBranch. Override any existing value with
// the srcItem's value. If the destination branch does not have an item
// with a matching keyName, append a copy of the srcItem.
func (dstBranch *Branch) mergeScalarItem(srcItem Item) {
	var dstItem *Item

	// if the destination already has an item with this key
	if dstBranch.ItemExists(srcItem.key) {
		item, _ := dstBranch.GetItem(srcItem.key)
		if item.Type() == "[leaf]" {
			item.value = srcItem.value
		}
		item.blockComments = srcItem.blockComments
		item.terminalWhitespace = srcItem.terminalWhitespace
		item.terminalComment = srcItem.terminalComment
		item.srcFile = srcItem.srcFile
		item.srcLine = srcItem.srcLine
		item.srcOrigin = srcItem.srcOrigin
		dstItem = item

	} else {
		// if the destination doesn't have an item with this key
		dupItem := srcItem.Copy()
		dstBranch.Items = append(dstBranch.Items, dupItem)
		dstItem = &dupItem
	}
	// recurse branches
	if srcItem.Type() == "[branch]" {
		innerDst := dstItem.value.(*Branch)
		innerSrc := srcItem.value.(*Branch)
		innerDst.Merge(innerSrc)
	}

}

// Merge the items with the given keyName. If items are only present in one of the two branches
// keep those items. If items are present in both branches, discard all of the dstBranch items
// and replace them with the srcBranch items.
// The dstBranch is typically a branch of the baselineTree
// The srcBranch is typically a branch of the userTree
func (dstBranch *Branch) mergeArrayItems(srcBranch *Branch, keyName string) {
	dstItems := dstBranch.QueryAll(keyName)
	srcItems := srcBranch.QueryAll(keyName)

	if len(srcItems) == 0 {
		return
	}
	if len(dstItems) == 0 {
		dstBranch.Items = append(dstBranch.Items, srcItems...)
		return
	}
	for {
		err := dstBranch.RemoveItem(keyName)
		if err == ErrNotFound {
			break
		}
	}
	dstBranch.Items = append(dstBranch.Items, srcItems...)
}
