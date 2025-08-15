package reviewer

import "github.com/spf13/cobra"

type Reviewers []Reviewer

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (reviewers Reviewers) GetHeaders(cmd *cobra.Command) []string {
	return Reviewer{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (reviewers Reviewers) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(reviewers) {
		return []string{}
	}
	return reviewers[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (reviewers Reviewers) Size() int {
	return len(reviewers)
}
