package branch

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Branch struct {
	Type                 string        `json:"type"                             mapstructure:"type"`
	Name                 string        `json:"name"                             mapstructure:"name"`
	Target               commit.Commit `json:"target"                           mapstructure:"target"`
	Links                link.Links    `json:"links"                            mapstructure:"links"`
	MergeStrategies      []string      `json:"merge_strategies,omitempty"       mapstructure:"merge_strategies"`
	DefaultMergeStrategy string        `json:"default_merge_strategy,omitempty" mapstructure:"default_merge_strategy"`
}

// Log is the logger for this application
var Log *logger.Logger

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "branch",
	Short: "Manage branches",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Branch requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// Validate validates a Branch
func (branch *Branch) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this Branch
//
// implements fmt.Stringer
func (branch Branch) String() string {
	return branch.Name
}
