package pullrequest

import "github.com/spf13/cobra"

type PullRequests []PullRequest

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (pullrequests PullRequests) GetHeaders(cmd *cobra.Command) []string {
	return PullRequest{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (pullrequests PullRequests) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(pullrequests) {
		return []string{}
	}
	return pullrequests[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (pullrequests PullRequests) Size() int {
	return len(pullrequests)
}
