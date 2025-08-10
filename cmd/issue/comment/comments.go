package comment

import "github.com/spf13/cobra"

type Comments []Comment

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (comments Comments) GetHeaders(cmd *cobra.Command) []string {
	return Comment{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Comments) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Comments) Size() int {
	return len(comments)
}
