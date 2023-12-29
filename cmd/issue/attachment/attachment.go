package attachment

import (
	"context"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Attachment struct {
	Type string         `json:"type"  mapstructure:"type"`
	Name string         `json:"name"  mapstructure:"name"`
	Link AttachmentLink `json:"links" mapstructure:"links"`
}

// Command represents this folder's command.
var Command = &cobra.Command{
	Use:   "attachment",
	Short: "Manage attachments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Attachment requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (attachment Attachment) GetHeader(short bool) []string {
	return []string{"Name", "URL"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (attachment Attachment) GetRow(headers []string) []string {
	return []string{
		attachment.Name,
		attachment.Link.String(),
	}
}

// Validate validates a Comment
func (attachment *Attachment) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (attachment Attachment) String() string {
	return attachment.Name
}

// GetIssueIDs gets the IDs of the issues
func GetIssueIDs(context context.Context, cmd *cobra.Command, args []string) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	type Issue struct {
		ID int `json:"id" mapstructure:"id"`
	}

	log.Infof("Getting all issues")
	issues, err := profile.GetAll[Issue](context, cmd, profile.Current, "issues")
	if err != nil {
		log.Errorf("Failed to get issues", err)
		return []string{}
	}
	return core.Map(issues, func(issue Issue) string {
		return fmt.Sprintf("%d", issue.ID)
	})
}

// GetAttachmentNames gets the names of the attachments
func GetAttachmentNames(context context.Context, cmd *cobra.Command, currentProfile *profile.Profile, issueID string) (names []string) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	log.Infof("Getting all attachments")
	attachments, err := profile.GetAll[Attachment](context, cmd, currentProfile, fmt.Sprintf("issues/%s/attachments", issueID))
	if err != nil {
		log.Errorf("Failed to get attachments", err)
		return []string{}
	}
	return core.Map(attachments, func(attachment Attachment) string {
		return attachment.Name
	})
}
