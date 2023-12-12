package workspace

type Workspaces []Workspace

// GetHeader gets the header for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetHeader() []string {
	return Workspace{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetRowAt(index int) []string {
	if index < 0 || index >= len(workspaces) {
		return []string{}
	}
	return workspaces[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (pullrequests Workspaces) Size() int {
	return len(pullrequests)
}
