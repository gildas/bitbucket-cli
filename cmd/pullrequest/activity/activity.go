package activity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/comment"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
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
	PullRequest PullRequestReference `json:"pull_request" mapstructure:"pullrequest"`
	Approval    *Approval            `json:"approval,omitempty"     mapstructure:"approval"`
	Comment     *comment.Comment     `json:"comment,omitempty"      mapstructure:"comment"`
	Update      *Update              `json:"update,omitempty"       mapstructure:"update"`
}

type Approval struct {
	Date        time.Time             `json:"date"        mapstructure:"date"`
	User        user.User             `json:"user"        mapstructure:"user"`
	PullRequest *PullRequestReference `json:"pullrequest" mapstructure:"pullrequest"`
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

var columns = common.Columns[Activity]{
	{Name: "pull_request", DefaultSorter: true, Compare: func(a, b Activity) bool {
		return a.PullRequest.ID < b.PullRequest.ID
	}},
	{Name: "date", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Approval != nil && b.Approval != nil {
			return a.Approval.Date.Before(b.Approval.Date)
		} else if a.Update != nil && b.Update != nil {
			return a.Update.Date.Before(b.Update.Date)
		}
		return false
	}},
	{Name: "approved", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Approval != nil && b.Approval != nil {
			return a.Approval.User.Name < b.Approval.User.Name
		}
		return false
	}},
	{Name: "description", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.Description), strings.ToLower(b.Update.Description)) == -1
		}
		return false
	}},
	{Name: "state", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.State), strings.ToLower(b.Update.State)) == -1
		}
		return false
	}},
	{Name: "author", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.Author.Name), strings.ToLower(b.Update.Author.Name)) == -1
		}
		return false
	}},
	{Name: "closed_by", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.ClosedBy.Name), strings.ToLower(b.Update.ClosedBy.Name)) == -1
		}
		return false
	}},
	{Name: "reason", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.Reason), strings.ToLower(b.Update.Reason)) == -1
		}
		return false
	}},
	{Name: "user", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Approval != nil && b.Approval != nil {
			return strings.Compare(strings.ToLower(a.Approval.User.Name), strings.ToLower(b.Approval.User.Name)) == -1
		} else if a.Update != nil && b.Update != nil {
			return strings.Compare(strings.ToLower(a.Update.Author.Name), strings.ToLower(b.Update.Author.Name)) == -1
		}
		return false
	}},
	{Name: "destination", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil && a.Update.Destination.Repository != nil && b.Update.Destination.Repository != nil {
			return strings.Compare(strings.ToLower(a.Update.Destination.Repository.Name), strings.ToLower(b.Update.Destination.Repository.Name)) == -1
		}
		return false
	}},
	{Name: "source", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil && a.Update.Source.Repository != nil && b.Update.Source.Repository != nil {
			return strings.Compare(strings.ToLower(a.Update.Source.Repository.Name), strings.ToLower(b.Update.Source.Repository.Name)) == -1
		}
		return false
	}},
	{Name: "created_on", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil {
			return a.Update.CreatedOn.Before(b.Update.CreatedOn)
		}
		return false
	}},
	{Name: "updated_on", DefaultSorter: false, Compare: func(a, b Activity) bool {
		if a.Update != nil && b.Update != nil && !a.Update.UpdatedOn.IsZero() && !b.Update.UpdatedOn.IsZero() {
			return a.Update.UpdatedOn.Before(b.Update.UpdatedOn)
		}
		return false
	}},
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (activity Activity) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Date", "Approved", "State", "User"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (activity Activity) GetRow(headers []string) []string {
	var row []string
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

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "date":
			row = append(row, activityDate.Format("2006-01-02 15:04:05"))
		case "approved":
			row = append(row, strconv.FormatBool(approval))
		case "description":
			if activity.Update != nil {
				row = append(row, activity.Update.Description)
			} else {
				row = append(row, " ")
			}
		case "state":
			row = append(row, state)
		case "author":
			if activity.Update != nil {
				row = append(row, activity.Update.Author.Name)
			} else {
				row = append(row, " ")
			}
		case "closed by":
			if activity.Update != nil {
				row = append(row, activity.Update.ClosedBy.Name)
			} else {
				row = append(row, " ")
			}
		case "reason":
			if activity.Update != nil {
				row = append(row, activity.Update.Reason)
			} else {
				row = append(row, " ")
			}
		case "user":
			row = append(row, user.Name)
		case "destination":
			if activity.Update != nil && activity.Update.Destination.Repository != nil {
				row = append(row, activity.Update.Destination.Repository.Name)
			} else {
				row = append(row, " ")
			}
		case "source":
			if activity.Update != nil && activity.Update.Source.Repository != nil {
				row = append(row, activity.Update.Source.Repository.Name)
			} else {
				row = append(row, " ")
			}
		case "created on", "created_on", "created-on", "created":
			if activity.Update != nil {
				row = append(row, activity.Update.CreatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		case "updated on", "updated_on", "updated-on", "updated":
			if activity.Update != nil && !activity.Update.UpdatedOn.IsZero() {
				row = append(row, activity.Update.UpdatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
}

// Validate validates a Comment
func (activity *Activity) Validate() error {
	var merr errors.MultiError

	if activity.Approval == nil && activity.Comment == nil && activity.Update == nil {
		merr.Append(errors.ArgumentMissing.With("approval, comment, or update"))
	}

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (activity Activity) String() string {
	return activity.PullRequest.Title
}

// MarshalJSON implements the json.Marshaler interface.
//
// implements json.Marshaler
func (activity Activity) MarshalJSON() (data []byte, err error) {
	type surrogate Activity

	data, err = json.Marshal(struct {
		surrogate
	}{
		surrogate: surrogate(activity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// implements json.Unmarshaler
func (activity *Activity) UnmarshalJSON(data []byte) (err error) {
	type surrogate Activity

	var surrogateActivity surrogate
	if err = json.Unmarshal(data, &surrogateActivity); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}

	*activity = Activity(surrogateActivity)
	return errors.JSONUnmarshalError.Wrap(activity.Validate())
}
