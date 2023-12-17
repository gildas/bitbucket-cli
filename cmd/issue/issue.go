package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/component"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
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
	Reporter   user.Account          `json:"reporter"   mapstructure:"reporter"`
	Assignee   user.Account          `json:"assignee"   mapstructure:"assignee"`
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

// GetHeader gets the header for a table
//
// implements common.Tableable
func (issue Issue) GetHeader(short bool) []string {
	return []string{"ID", "Title", "State", "Priority", "Repository", "Reporter", "Assignee"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (issue Issue) GetRow(headers []string) []string {
	return []string{
		fmt.Sprintf("%d", issue.ID),
		issue.Title,
		issue.State,
		issue.Priority,
		issue.Repository.Name,
		issue.Reporter.Name,
		issue.Assignee.Name,
	}
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
func GetIssueIDs(context context.Context, p *profile.Profile, repository string) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	issues, err := profile.GetAll[Issue](
		context,
		p,
		repository,
		"issues",
	)
	if err != nil {
		log.Errorf("Failed to get issues", err)
		return []string{}
	}
	ids = make([]string, 0, len(issues))
	for _, issue := range issues {
		ids = append(ids, fmt.Sprintf("%d", issue.ID))
	}
	return
}
