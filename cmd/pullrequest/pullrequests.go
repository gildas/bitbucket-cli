package pullrequest

import (
	"fmt"
	"strings"
)

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

// PullRequestsWithOptions is a wrapper that supports optional columns
type PullRequestsWithOptions struct {
	pullrequests    []PullRequest
	showDescription bool
	showCreatedAt   bool
}

// NewPullRequestsWithOptions creates a new PullRequestsWithOptions
func NewPullRequestsWithOptions(pullrequests []PullRequest, showDescription, showCreatedAt bool) PullRequestsWithOptions {
	return PullRequestsWithOptions{
		pullrequests:    pullrequests,
		showDescription: showDescription,
		showCreatedAt:   showCreatedAt,
	}
}

// GetHeader gets the header for a table
//
// implements common.Tableables
func (p PullRequestsWithOptions) GetHeader() []string {
	headers := []string{"ID", "TITLE"}
	if p.showDescription {
		headers = append(headers, "DESCRIPTION")
	}
	headers = append(headers, "BRANCH")
	if p.showCreatedAt {
		headers = append(headers, "CREATED AT")
	}
	headers = append(headers, "STATE")
	return headers
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (p PullRequestsWithOptions) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(p.pullrequests) {
		return []string{}
	}
	pullrequest := p.pullrequests[index]

	row := []string{
		fmt.Sprintf("#%d", pullrequest.ID),
		pullrequest.Title,
	}

	if p.showDescription {
		row = append(row, pullrequest.Description)
	}

	row = append(row, pullrequest.Source.Branch.Name)

	if p.showCreatedAt {
		row = append(row, pullrequest.CreatedOn.Format("2006-01-02 15:04:05"))
	}

	row = append(row, strings.ToUpper(pullrequest.State))
	return row
}

// Size gets the number of elements
//
// implements common.Tableables
func (p PullRequestsWithOptions) Size() int {
	return len(p.pullrequests)
}
