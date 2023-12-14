package reviewer

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Reviewer struct {
	Type         string    `json:"type" mapstructure:"type"`
	ReviewerType string    `json:"reviewer_type" mapstructure:"reviewer_type"`
	User         user.User `json:"user" mapstructure:"user"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "reviewer",
	Short: "Manage reviewers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reviewer requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetHeader(short bool) []string {
	return []string{"Type", "Reviewer Type", "User"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetRow(headers []string) []string {
	return []string{reviewer.Type, reviewer.ReviewerType, reviewer.User.Name}
}

// Validate validates a Reviewer
func (reviewer *Reviewer) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}
