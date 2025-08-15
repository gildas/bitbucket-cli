package issue

import "github.com/spf13/cobra"

type Issues []Issue

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (issues Issues) GetHeaders(cmd *cobra.Command) []string {
	return Issue{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (issues Issues) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(issues) {
		return []string{}
	}
	return issues[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (issues Issues) Size() int {
	return len(issues)
}
