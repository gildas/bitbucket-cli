package commit

import "github.com/spf13/cobra"

type Commits []Commit

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (commits Commits) GetHeaders(cmd *cobra.Command) []string {
	return Commit{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (commits Commits) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(commits) {
		return []string{}
	}
	return commits[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (commits Commits) Size() int {
	return len(commits)
}
