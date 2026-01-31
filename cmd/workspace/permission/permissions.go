package permission

import "github.com/spf13/cobra"

type Permissions []Permission

// GetHeaders get the headers for a table
//
// implements common.Tableables
func (permissions Permissions) GetHeaders(cmd *cobra.Command) []string {
	return Permission{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (permissions Permissions) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(permissions) {
		return []string{}
	}
	return permissions[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (permissions Permissions) Size() int {
	return len(permissions)
}
