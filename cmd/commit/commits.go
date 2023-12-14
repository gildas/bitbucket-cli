package commit

type Commits []Commit

// GetHeader gets the header for a table
//
// implements common.Tableables
func (commits Commits) GetHeader() []string {
	return Commit{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (commits Commits) GetRowAt(index int) []string {
	if index < 0 || index >= len(commits) {
		return []string{}
	}
	return commits[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (commits Commits) Size() int {
	return len(commits)
}
