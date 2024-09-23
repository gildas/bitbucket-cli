package pullrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/comment"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PullRequest struct {
	Type              string              `json:"type"                   mapstructure:"type"`
	ID                uint64              `json:"id"                     mapstructure:"id"`
	Title             string              `json:"title"                  mapstructure:"title"`
	Description       string              `json:"description"            mapstructure:"description"`
	Summary           common.RenderedText `json:"summary"                mapstructure:"summary"`
	State             string              `json:"state"                  mapstructure:"state"`
	MergeCommit       *commit.Commit      `json:"merge_commit,omitempty" mapstructure:"merge_commit"`
	CloseSourceBranch bool                `json:"close_source_branch"    mapstructure:"close_source_branch"`
	ClosedBy          user.User           `json:"closed_by"              mapstructure:"closed_by"`
	Author            user.Account        `json:"author"                 mapstructure:"author"`
	Reason            string              `json:"reason"                 mapstructure:"reason"`
	Destination       Endpoint            `json:"destination"            mapstructure:"destination"`
	Source            Endpoint            `json:"source"                 mapstructure:"source"`
	Links             common.Links        `json:"links"                  mapstructure:"links"`
	CommentCount      uint64              `json:"comment_count"          mapstructure:"comment_count"`
	TaskCount         uint64              `json:"task_count"             mapstructure:"task_count"`
	CreatedOn         time.Time           `json:"created_on"             mapstructure:"created_on"`
	UpdatedOn         time.Time           `json:"updated_on"             mapstructure:"updated_on"`
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

func init() {
	Command.AddCommand(comment.Command)
}

func SetLogger(log *logger.Logger) {
	createOptions.Reviewers = createOptions.Reviewers.WithLogger(log)
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (pullrequest PullRequest) GetHeader(short bool) []string {
	return []string{"ID", "Title", "Description", "source", "destination", "state"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (pullrequest PullRequest) GetRow(headers []string) []string {
	return []string{
		fmt.Sprintf("%d", pullrequest.ID),
		pullrequest.Title,
		pullrequest.Description,
		pullrequest.Source.Branch.Name,
		pullrequest.Destination.Branch.Name,
		pullrequest.State,
	}
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

// GetPullRequestIDs gets the pullrequest Ids for completion
func GetPullRequestIDs(context context.Context, cmd *cobra.Command, repository string, state string) []string {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getpullrequests")

	log.Infof("Getting open pullrequests for repository %s", approveOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(context),
		cmd,
		profile.Current,
		fmt.Sprintf("pullrequests?state=%s", state),
	)
	if err != nil {
		log.Errorf("Failed to get pullrequests for repository %s", repository, err)
		return []string{}
	}

	return core.Map(pullrequests, func(pullrequest PullRequest) string {
		return fmt.Sprintf("%d", pullrequest.ID)
	})
}

// GetReviewerNicknames gets the reviewer nicknames for the current Workspace
func GetReviewerNicknames(context context.Context, cmd *cobra.Command, args []string) []string {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getreviewers")
	var pullrequestWorkspace *workspace.Workspace
	var err error

	if workspaceName := cmd.Flag("workspace").Value.String(); len(workspaceName) > 0 {
		pullrequestWorkspace, err = workspace.GetWorkspace(cmd.Context(), cmd, profile.Current, workspaceName)
	} else {
		pullrequestWorkspace, err = workspace.GetWorkspaceFromGit(cmd.Context(), cmd, profile.Current)
	}
	if err != nil {
		log.Errorf("Failed to get repository: %s", err)
		return []string{}
	}
	members, _ := pullrequestWorkspace.GetMembers(context, cmd)
	return core.Map(members, func(member workspace.Member) string {
		return member.User.Nickname
	})
}

