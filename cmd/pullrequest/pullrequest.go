package pullrequest

import (
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
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
	MergeCommit       *commit.Commit     `json:"merge_commit,omitempty" mapstructure:"merge_commit"`
	CloseSourceBranch bool               `json:"close_source_branch"    mapstructure:"close_source_branch"`
	ClosedBy          user.User          `json:"closed_by"              mapstructure:"closed_by"`
	Author            common.AppUser     `json:"author"                 mapstructure:"author"`
	Reason            string             `json:"reason"                 mapstructure:"reason"`
	Destination       Endpoint           `json:"destination"            mapstructure:"destination"`
	Source            Endpoint           `json:"source"                 mapstructure:"source"`
	Links             link.Links         `json:"links"                  mapstructure:"links"`
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

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr", "pull-request"},
	Short:   "Manage pull requests",
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

// getOpenPullRequests gets the pullrequests for completion
func getOpenPullRequests(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")

	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Infof("Getting open pullrequests for repository %s", approveOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"pullrequests?state=OPEN",
	)
	if err != nil {
		log.Errorf("Failed to get pullrequests for repository %s", unapproveOptions.Repository, err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var result []string
	for _, pullrequest := range pullrequests {
		result = append(result, fmt.Sprintf("%d", pullrequest.ID))
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}
