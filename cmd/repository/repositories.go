package repository

import "github.com/spf13/cobra"

type Repositories []Repository

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (repositories Repositories) GetHeaders(cmd *cobra.Command) []string {
	return Repository{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (repositories Repositories) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(repositories) {
		return []string{}
	}
	return repositories[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (repositories Repositories) Size() int {
	return len(repositories)
}
