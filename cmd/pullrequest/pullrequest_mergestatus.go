package pullrequest

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

// PullRequestMergeStatus describes the status of a pull request
type PullRequestMergeStatus struct {
	ID          string       `json:"id"          mapstructure:"id"`
	Status      string       `json:"task_status" mapstructure:"task_status"`
	PullRequest PullRequest  `json:"merge_result" mapstructure:"merge_result"`
	Links       common.Links `json:"links"       mapstructure:"links"`
}

// NewPullRequestMergeStatusFromLocation creates a new PullRequestMergeStatus from a URL location
//
// The URL location is the URL returned in the Location header of the response from Bitbucket when we request to merge a pull request asynchronously
func NewPullRequestMergeStatusFromLocation(location string) (mergeStatus *PullRequestMergeStatus, err error) {
	// Format: https://api.bitbucket.org/2.0/repositories/<workspace_slug>/<repo_slug>/pullrequests/<pullrequest_id>/merge/task-status/<task_id>
	if len(location) == 0 {
		return nil, errors.Errorf("Failed to get the merge task URL from the Location header in the response from Bitbucket")
	}
	parts := strings.Split(location, "/")
	if len(parts) < 2 {
		return nil, errors.Errorf("Invalid merge task URL: %s", location)
	}
	taskID := parts[len(parts)-1]
	pullrequestID, err := strconv.Atoi(parts[len(parts)-5])
	if err != nil {
		return nil, errors.Errorf("Invalid pull request ID: %s", parts[len(parts)-5])
	}
	return &PullRequestMergeStatus{ID: taskID, PullRequest: PullRequest{ID: uint64(pullrequestID)}}, nil
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (status PullRequestMergeStatus) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Pull Request", "Status"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (status PullRequestMergeStatus) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, status.ID)
		case "pull request", "pull_request", "pull-request", "pullrequest", "pr":
			row = append(row, fmt.Sprintf("%d", status.PullRequest.ID))
		case "status":
			if status.Status == "SUCCESS" {
				row = append(row, status.Status+"/"+status.PullRequest.State)
			} else {
				row = append(row, status.Status)
			}
		}
	}
	return row
}
