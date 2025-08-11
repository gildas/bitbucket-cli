package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/component"
	"bitbucket.org/gildas_cherruel/bb/cmd/issue/attachment"
	"bitbucket.org/gildas_cherruel/bb/cmd/issue/comment"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Issue struct {
	Type       string                `json:"type"       mapstructure:"type"`
	ID         int                   `json:"id"         mapstructure:"id"`
	Kind       string                `json:"kind"       mapstructure:"kind"` // bug, enhancement, proposal, task
	Title      string                `json:"title"      mapstructure:"title"`
	State      string                `json:"state"      mapstructure:"state"`    // new, open, submitted, resolved, on hold, invalid, duplicate, wontfix, closed
	Priority   string                `json:"priority"   mapstructure:"priority"` // trivial, minor, major, critical, blocker
	Repository repository.Repository `json:"repository" mapstructure:"repository"`
	Reporter   user.User             `json:"reporter"   mapstructure:"reporter"`
	Assignee   user.User             `json:"assignee"   mapstructure:"assignee"`
	Content    common.RenderedText   `json:"content"    mapstructure:"content"`
	Votes      int                   `json:"votes"      mapstructure:"votes"`
	Watchers   int                   `json:"watches"    mapstructure:"watches"`
	Milestone  *common.Entity        `json:"milestone"  mapstructure:"milestone"`
	Component  *component.Component  `json:"component"  mapstructure:"component"`
	Links      common.Links          `json:"links"      mapstructure:"links"`
	CreatedOn  time.Time             `json:"created_on" mapstructure:"created_on"`
	UpdatedOn  time.Time             `json:"updated_on" mapstructure:"updated_on"`
	EditedOn   time.Time             `json:"edited_on"  mapstructure:"edited_on"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Issue requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = []string{
	"id",
	"title",
	"state",
	"priority",
	"repository",
	"reporter",
	"assignee",
	"created_on",
	"updated_on",
	"edited_on",
	"votes",
	"watchers",
	"milestone",
}

var sortBy = []string{
	"+id",
	"title",
	"state",
	"priority",
	"repository",
	"reporter",
	"assignee",
	"created_on",
	"updated_on",
	"edited_on",
	"votes",
	"watchers",
	"milestone",
}

func init() {
	Command.AddCommand(comment.Command)
	Command.AddCommand(attachment.Command)
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (issue Issue) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Title", "State", "Priority", "Repository", "Reporter", "Assignee"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (issue Issue) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, fmt.Sprintf("%d", issue.ID))
		case "kind":
			row = append(row, issue.Kind)
		case "title":
			row = append(row, issue.Title)
		case "state":
			row = append(row, issue.State)
		case "priority":
			row = append(row, issue.Priority)
		case "repository":
			row = append(row, issue.Repository.Name)
		case "reporter":
			row = append(row, issue.Reporter.Name)
		case "assignee":
			row = append(row, issue.Assignee.Name)
		case "created on", "created_on", "created-on":
			row = append(row, issue.CreatedOn.Format("2006-01-02 15:04:05"))
		case "updated on", "updated_on", "updated-on":
			if !issue.UpdatedOn.IsZero() {
				row = append(row, issue.UpdatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		case "edited on", "edited_on", "edited-on":
			if !issue.EditedOn.IsZero() {
				row = append(row, issue.EditedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		case "votes":
			row = append(row, fmt.Sprintf("%d", issue.Votes))
		case "watchers":
			row = append(row, fmt.Sprintf("%d", issue.Watchers))
		case "milestone":
			if issue.Milestone != nil {
				row = append(row, issue.Milestone.Name)
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
}

// Validate validates a Issue
func (issue *Issue) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (issue Issue) String() string {
	return issue.Title
}

// MarshalJSON implements the json.Marshaler interface.
func (issue Issue) MarshalJSON() (data []byte, err error) {
	type surrogate Issue

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
		EditedOn  string `json:"edited_on"`
	}{
		surrogate: surrogate(issue),
		CreatedOn: issue.CreatedOn.Format(time.RFC3339),
		UpdatedOn: issue.UpdatedOn.Format(time.RFC3339),
		EditedOn:  issue.EditedOn.Format(time.RFC3339),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// GetIssueIDs gets the IDs of the issues
func GetIssueIDs(context context.Context, cmd *cobra.Command) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	issues, err := profile.GetAll[Issue](context, cmd, "issues")
	if err != nil {
		log.Errorf("Failed to get issues", err)
		return
	}
	ids = core.Map(issues, func(issue Issue) string { return fmt.Sprintf("%d", issue.ID) })
	core.Sort(ids, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return ids, nil
}
