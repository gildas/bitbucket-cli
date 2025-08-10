package workspace

import "github.com/spf13/cobra"

type Workspaces []Workspace

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetHeaders(cmd *cobra.Command) []string {
	return Workspace{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(workspaces) {
		return []string{}
	}
	return workspaces[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (pullrequests Workspaces) Size() int {
	return len(pullrequests)
}
