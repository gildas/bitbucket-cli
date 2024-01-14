package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Comment struct {
	Type        string                `json:"type"             mapstructure:"type"`
	ID          int                   `json:"id"               mapstructure:"id"`
	Content     common.RenderedText   `json:"content"          mapstructure:"content"`
	User        user.Account          `json:"user"             mapstructure:"user"`
	Anchor      *common.FileAnchor    `json:"inline,omitempty" mapstructure:"inline"`
	Parent      *Comment              `json:"parent,omitempty" mapstructure:"parent"`
	CreatedOn   time.Time             `json:"created_on"       mapstructure:"created_on"`
	UpdatedOn   time.Time             `json:"updated_on"       mapstructure:"updated_on"`
	IsDeleted   bool                  `json:"deleted"          mapstructure:"deleted"`
	IsPending   bool                  `json:"pending"          mapstructure:"pending"`
	PullRequest *PullRequestReference `json:"pullrequest"      mapstructure:"pullrequest"`
	Links       common.Links          `json:"links"            mapstructure:"links"`
}

type PullRequestReference struct {
	Type  string       `json:"type"  mapstructure:"type"`
	ID    int          `json:"id"    mapstructure:"id"`
	Title string       `json:"title" mapstructure:"title"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "comment",
	Short: "Manage comments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Comment requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (comment Comment) GetHeader(short bool) []string {
	if short {
		headers := []string{"ID", "Created On"}
		if !comment.UpdatedOn.IsZero() {
			headers = append(headers, "Updated On")
		}
		if comment.Anchor != nil {
			headers = append(headers, "File")
		}
		return append(headers, "User", "Content")
	}
	return []string{"ID", "Created On", "Updated On", "File", "User", "Content"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (comment Comment) GetRow(headers []string) []string {
	rows := []string{
		fmt.Sprintf("%d", comment.ID),
		comment.CreatedOn.Format("2006-01-02 15:04:05"),
	}
	if core.Contains(headers, "Updated On") {
		updatedOn := ""
		if !comment.UpdatedOn.IsZero() {
			updatedOn = comment.UpdatedOn.Format("2006-01-02 15:04:05")
		}
		rows = append(rows, updatedOn)
	}

	if core.Contains(headers, "File") {
		file := ""
		if comment.Anchor != nil {
			file = comment.Anchor.String()
		}
		rows = append(rows, file)
	}

	return append(rows,
		comment.User.Name,
		comment.Content.Raw,
	)
}

// Validate validates a Comment
func (comment *Comment) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (comment Comment) String() string {
	return comment.Content.Raw
}

// MarshalJSON implements the json.Marshaler interface.
func (comment Comment) MarshalJSON() (data []byte, err error) {
	type surrogate Comment

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
	}{
		surrogate: surrogate(comment),
		CreatedOn: comment.CreatedOn.Format(time.RFC3339),
		UpdatedOn: comment.UpdatedOn.Format(time.RFC3339),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// GetPullRequestIDs gets the IDs of the pullrequests
func GetPullRequestIDs(context context.Context, cmd *cobra.Command, args []string) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("pullrequest", "getids")

	type PullRequest struct {
		ID int `json:"id" mapstructure:"id"`
	}

	log.Infof("Getting all pullrequests")
	pullrequests, err := profile.GetAll[PullRequest](context, cmd, profile.Current, "pullrequests")
	if err != nil {
		log.Errorf("Failed to get pullrequests", err)
		return []string{}
	}
	return core.Map(pullrequests, func(pullrequest PullRequest) string {
		return fmt.Sprintf("%d", pullrequest.ID)
	})
}

// GetPullRequestCommentIDs gets the IDs of the comments for a pullrequest
func GetPullRequestCommentIDs(context context.Context, cmd *cobra.Command, currentProfile *profile.Profile, PullRequestID string) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("pullrequest", "getids")

	comments, err := profile.GetAll[Comment](context, cmd, currentProfile, fmt.Sprintf("pullrequests/%s/comments", PullRequestID))
	if err != nil {
		log.Errorf("Failed to get pullrequests", err)
		return []string{}
	}
	return core.Map(comments, func(comment Comment) string {
		return fmt.Sprintf("%d", comment.ID)
	})
}
