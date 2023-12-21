package pullrequest

type PullRequests []PullRequest

// GetHeader gets the header for a table
//
// implements common.Tableables
func (pullrequests PullRequests) GetHeader() []string {
	return PullRequest{}.GetHeader(false)
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
