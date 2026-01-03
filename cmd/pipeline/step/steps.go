package step

import "github.com/spf13/cobra"

type Steps []Step

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (comments Steps) GetHeaders(cmd *cobra.Command) []string {
	return Step{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Steps) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Steps) Size() int {
	return len(comments)
}
