package pullrequest

import (
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PullRequest struct {
	Type              string             `json:"type"                   mapstructure:"type"`
	ID                uint64             `json:"id"                     mapstructure:"id"`
	Title             string             `json:"title"                  mapstructure:"title"`
	Description       string             `json:"description"            mapstructure:"description"`
	Summary           PullRequestSummary `json:"summary"                mapstructure:"summary"`
	State             string             `json:"state"                  mapstructure:"state"`
	MergeCommit       *Commit            `json:"merge_commit,omitempty" mapstructure:"merge_commit"`
	CloseSourceBranch bool               `json:"close_source_branch"    mapstructure:"close_source_branch"`
	ClosedBy          common.User        `json:"closed_by"              mapstructure:"closed_by"`
	Author            common.AppUser     `json:"author"                 mapstructure:"author"`
	Reason            string             `json:"reason"                 mapstructure:"reason"`
	Destination       Endpoint           `json:"destination"            mapstructure:"destination"`
	Source            Endpoint           `json:"source"                 mapstructure:"source"`
	Links             common.Links       `json:"links"                  mapstructure:"links"`
	CommentCount      uint64             `json:"comment_count"          mapstructure:"comment_count"`
	TaskCount         uint64             `json:"task_count"             mapstructure:"task_count"`
	CreatedOn         time.Time          `json:"created_on"             mapstructure:"created_on"`
	UpdatedOn         time.Time          `json:"updated_on"             mapstructure:"updated_on"`
}

type PullRequestSummary struct {
	Type   string `json:"type"   mapstructure:"type"`
	Markup string `json:"markup" mapstructure:"markup"`
	Raw    string `json:"raw"    mapstructure:"raw"`
	HTML   string `json:"html"   mapstructure:"html"`
}

// Log is the logger for this application
var Log *logger.Logger

// Profile is the profile for this command
var Profile *profile.Profile

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "pullrequest",
	Short: "Manage pull requests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pullrequest requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// Validate validates a PullRequest
func (pullrequest *PullRequest) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (pullrequest PullRequest) String() string {
	return pullrequest.Title
}

// MarshalJSON implements the json.Marshaler interface.
func (pullrequest PullRequest) MarshalJSON() (data []byte, err error) {
	type surrogate PullRequest

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
	}{
		surrogate: surrogate(pullrequest),
		CreatedOn: pullrequest.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
		UpdatedOn: pullrequest.UpdatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
