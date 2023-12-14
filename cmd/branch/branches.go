package branch

type Branches []Branch

// GetHeader gets the header for a table
//
// implements common.Tableables
func (branches Branches) GetHeader() []string {
	return Branch{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (branches Branches) GetRowAt(index int) []string {
	if index < 0 || index >= len(branches) {
		return []string{}
	}
	return branches[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (branches Branches) Size() int {
	return len(branches)
}
