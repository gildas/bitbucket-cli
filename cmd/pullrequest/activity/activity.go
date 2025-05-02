package activity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Branch struct {
	Name                 string   `json:"name"                             mapstructure:"name"`
	MergeStrategies      []string `json:"merge_strategies,omitempty"       mapstructure:"merge_strategies"`
	DefaultMergeStrategy string   `json:"default_merge_strategy,omitempty" mapstructure:"default_merge_strategy"`
}

type Endpoint struct {
	Branch     Branch                 `json:"branch"               mapstructure:"branch"`
	Commit     *commit.Commit         `json:"commit,omitempty"     mapstructure:"commit"`
	Repository *repository.Repository `json:"repository,omitempty" mapstructure:"repository"`
}

type Activity struct {
	PullRequest *PullRequestReference `json:"pull_request" mapstructure:"pullrequest"`
	Approval    *Approval             `json:"approval,omitempty"     mapstructure:"approval"`
	Update      *Update               `json:"update,omitempty"       mapstructure:"update"`
}

type Update struct {
	Date              time.Time           `json:"date"                   mapstructure:"date"`
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

type Approval struct {
	Date        time.Time             `json:"date"        mapstructure:"date"`
	User        user.User             `json:"user"        mapstructure:"user"`
	PullRequest *PullRequestReference `json:"pullrequest" mapstructure:"pullrequest"`
}

type PullRequestReference struct {
	Type  string       `json:"type"  mapstructure:"type"`
	ID    int          `json:"id"    mapstructure:"id"`
	Title string       `json:"title" mapstructure:"title"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "activity",
	Short: "Manage activities",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Activity requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (activity Activity) GetHeader(short bool) []string {
	return []string{"Date", "Approved", "State", "User"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (activity Activity) GetRow(headers []string) []string {
	rows := []string{}

	var activityDate time.Time
	var approval bool
	var state string
	var user user.User

	if activity.Approval != nil {
		activityDate = activity.Approval.Date
		approval = true
		user = activity.Approval.User
		state = "N/A"
	} else if activity.Update != nil {
		activityDate = activity.Update.Date
		state = activity.Update.State
		user = activity.Update.Author
		approval = false
	}

	return append(rows,
		activityDate.Format("2006-01-02 15:04:05"),
		strconv.FormatBool(approval),
		state,
		user.Name,
	)
}

// Validate validates a Comment
func (activity *Activity) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (activity Activity) String() string {
	return activity.PullRequest.Title
}

// MarshalJSON implements the json.Marshaler interface.
func (activity Activity) MarshalJSON() (data []byte, err error) {
	type surrogate Activity

	data, err = json.Marshal(struct {
		surrogate
	}{
		surrogate: surrogate(activity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
