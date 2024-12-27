package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/spf13/cobra"
)

type IssueChange struct {
	Type      string              `json:"type"       mapstructure:"type"`
	ID        int                 `json:"id"         mapstructure:"id"`
	Name      string              `json:"name"       mapstructure:"name"`
	Issue     Issue               `json:"issue"      mapstructure:"issue"`
	User      user.User           `json:"user"       mapstructure:"user"`
	Changes   ChangeSet           `json:"changes"    mapstructure:"changes"`
	Message   common.RenderedText `json:"message"    mapstructure:"message"`
	CreatedOn time.Time           `json:"created_on" mapstructure:"created_on"`
	Links     common.Links        `json:"links"      mapstructure:"links"`
}

type ChangeSet struct {
	Assignee  *Change `json:"assignee,omitempty"  mapstructure:"assignee,omitempty"`
	Kind      *Change `json:"kind,omitempty"      mapstructure:"kind,omitempty"`
	Priority  *Change `json:"priority,omitempty"  mapstructure:"priority,omitempty"`
	State     *Change `json:"state,omitempty"     mapstructure:"state,omitempty"`
	Title     *Change `json:"title,omitempty"     mapstructure:"title,omitempty"`
	Content   *Change `json:"content,omitempty"   mapstructure:"content,omitempty"`
	Milestone *Change `json:"milestone,omitempty" mapstructure:"milestone,omitempty"`
	Component *Change `json:"component,omitempty" mapstructure:"component,omitempty"`
	Version   *Change `json:"version,omitempty"   mapstructure:"version,omitempty"`
}

type Change struct {
	Old string `json:"old" mapstructure:"old,omitempty"`
	New string `json:"new" mapstructure:"new"`
}

// String gets a string representation
//
// implements fmt.Stringer
func (changeSet ChangeSet) String() string {
	var changes []string
	if changeSet.Assignee != nil {
		changes = append(changes, fmt.Sprintf("assignee: %s", changeSet.Assignee))
	}
	if changeSet.Kind != nil {
		changes = append(changes, fmt.Sprintf("kind: %s", changeSet.Kind))
	}
	if changeSet.Priority != nil {
		changes = append(changes, fmt.Sprintf("priority: %s", changeSet.Priority))
	}
	if changeSet.State != nil {
		changes = append(changes, fmt.Sprintf("state: %s", changeSet.State))
	}
	if changeSet.Title != nil {
		changes = append(changes, fmt.Sprintf("title: %s", changeSet.Title))
	}
	if changeSet.Content != nil {
		changes = append(changes, fmt.Sprintf("content: %s", changeSet.Content))
	}
	if changeSet.Milestone != nil {
		changes = append(changes, fmt.Sprintf("milestone: %s", changeSet.Milestone))
	}
	if changeSet.Component != nil {
		changes = append(changes, fmt.Sprintf("component: %s", changeSet.Component))
	}
	if changeSet.Version != nil {
		changes = append(changes, fmt.Sprintf("version: %s", changeSet.Version))
	}
	return strings.Join(changes, ", ")
}

// String gets a string representation
//
// implements fmt.Stringer
func (change Change) String() string {
	return fmt.Sprintf("%s -> %s", change.Old, change.New)
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (issueChange IssueChange) GetHeader(short bool) []string {
	return []string{"Issue ID", "Issue Title", "User", "Date", "Changes"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (issueChange IssueChange) GetRow(headers []string) []string {
	return []string{
		fmt.Sprintf("%d", issueChange.Issue.ID),
		issueChange.Issue.Title,
		issueChange.User.Name,
		issueChange.CreatedOn.Local().Format(time.RFC3339),
		issueChange.Changes.String(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (change IssueChange) MarshalJSON() (data []byte, err error) {
	type surrogate IssueChange

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
		EditedOn  string `json:"edited_on"`
	}{
		surrogate: surrogate(change),
		CreatedOn: change.CreatedOn.Format(time.RFC3339),
	})
	return
}

// GetIssueChanges gets the changes for an issue
func GetIssueChanges(context context.Context, cmd *cobra.Command, issueID string) (changes IssueChanges, err error) {
	return profile.GetAll[IssueChange](
		cmd.Context(),
		cmd,
		fmt.Sprintf("issues/%s/changes", issueID),
	)
}
