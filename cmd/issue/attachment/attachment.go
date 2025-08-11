package attachment

import (
	"context"
	"fmt"
	"strings"

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

var columns = []string{
	"name",
	"url",
	"type",
}

var sortBy = []string{
	"+name",
	"url",
	"type",
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (attachment Attachment) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Name", "URL"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (attachment Attachment) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "name":
			row = append(row, attachment.Name)
		case "link", "url":
			row = append(row, attachment.Link.String())
		case "type":
			row = append(row, attachment.Type)
		}
	}
	return row
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
func GetIssueIDs(context context.Context, cmd *cobra.Command, args []string, toComplete string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	type Issue struct {
		ID int `json:"id" mapstructure:"id"`
	}

	log.Infof("Getting all issues")
	issues, err := profile.GetAll[Issue](context, cmd, "issues")
	if err != nil {
		log.Errorf("Failed to get issues", err)
		return []string{}, err
	}
	ids = core.Map(issues, func(issue Issue) string { return fmt.Sprintf("%d", issue.ID) })
	core.Sort(ids, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return ids, nil
}

// GetAttachmentNames gets the names of the attachments
func GetAttachmentNames(context context.Context, cmd *cobra.Command, issueID string) (names []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("issue", "getids")

	log.Infof("Getting all attachments")
	attachments, err := profile.GetAll[Attachment](context, cmd, fmt.Sprintf("issues/%s/attachments", issueID))
	if err != nil {
		log.Errorf("Failed to get attachments", err)
		return []string{}, err
	}
	names = core.Map(attachments, func(attachment Attachment) string { return attachment.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return names, nil
}
