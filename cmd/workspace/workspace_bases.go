package workspace

import "github.com/spf13/cobra"

type WorkspaceBases []WorkspaceBase

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (workspaces WorkspaceBases) GetHeaders(cmd *cobra.Command) []string {
	return WorkspaceBase{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (workspaces WorkspaceBases) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(workspaces) {
		return []string{}
	}
	return workspaces[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (workspaces WorkspaceBases) Size() int {
	return len(workspaces)
}
