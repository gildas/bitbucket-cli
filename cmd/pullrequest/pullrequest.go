package pullrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/branch"
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/activity"
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
	Author            user.User           `json:"author"                 mapstructure:"author"`
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

var columns = []string{
	"id",
	"title",
	"description",
	"source",
	"destination",
	"state",
	"author",
	"closed_by",
	"commit",
	"reason",
	"comments",
	"tasks",
	"created_on",
	"updated_on",
}

func init() {
	Command.AddCommand(comment.Command)
	Command.AddCommand(activity.Command)
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (pullrequest PullRequest) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Title", "Description", "source", "destination", "state"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (pullrequest PullRequest) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, fmt.Sprintf("%d", pullrequest.ID))
		case "title":
			row = append(row, pullrequest.Title)
		case "description":
			row = append(row, pullrequest.Description)
		case "source":
			row = append(row, pullrequest.Source.Branch.Name)
		case "destination":
			row = append(row, pullrequest.Destination.Branch.Name)
		case "state":
			row = append(row, pullrequest.State)
		case "author":
			row = append(row, pullrequest.Author.Name)
		case "closed by":
			row = append(row, pullrequest.ClosedBy.Name)
		case "commit":
			if pullrequest.MergeCommit != nil {
				row = append(row, pullrequest.MergeCommit.Hash[:7])
			} else {
				row = append(row, " ")
			}
		case "reason":
			row = append(row, pullrequest.Reason)
		case "comments":
			row = append(row, fmt.Sprintf("%d", pullrequest.CommentCount))
		case "tasks":
			row = append(row, fmt.Sprintf("%d", pullrequest.TaskCount))
		case "created on", "created_on", "created-on":
			row = append(row, pullrequest.CreatedOn.Format("2006-01-02 15:04:05"))
		case "updated on", "updated_on", "updated-on":
			if !pullrequest.UpdatedOn.IsZero() {
				row = append(row, pullrequest.UpdatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
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

// GetReviewerNicknames gets the reviewer nicknames for the current Workspace
func GetReviewerNicknames(context context.Context, cmd *cobra.Command, args []string, toComplete string) (nicknames []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getreviewers")
	var pullrequestWorkspace *workspace.Workspace

	if cmd == nil {
		fmt.Fprintln(os.Stderr, "cmd is nil")
		return []string{}, errors.ArgumentMissing.With("cmd")
	}

	if workspaceName := cmd.Flag("workspace").Value.String(); len(workspaceName) > 0 {
		pullrequestWorkspace, err = workspace.GetWorkspace(cmd.Context(), cmd, workspaceName)
	} else {
		pullrequestWorkspace, err = workspace.GetWorkspaceFromGit(cmd.Context(), cmd)
	}
	if err != nil {
		log.Errorf("Failed to get repository: %s", err)
		return []string{}, err
	}
	members, _ := pullrequestWorkspace.GetMembers(context, cmd)
	nicknames = core.Map(members, func(member workspace.Member) string { return member.User.Nickname })
	core.Sort(nicknames, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return nicknames, nil
}

// GetBranchNames gets the branch names of a repository
func GetBranchNames(context context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getbranches")
	log.Infof("Getting branches for profile %v", profile.Current)
	names, err := branch.GetBranchNames(context, cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, err
	}
	return common.FilterValidArgs(names, args, toComplete), nil
}
