package branch

import (
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Branch struct {
	Type                 string        `json:"type"                             mapstructure:"type"`
	Name                 string        `json:"name"                             mapstructure:"name"`
	Target               commit.Commit `json:"target"                           mapstructure:"target"`
	Links                common.Links  `json:"links"                            mapstructure:"links"`
	MergeStrategies      []string      `json:"merge_strategies,omitempty"       mapstructure:"merge_strategies"`
	DefaultMergeStrategy string        `json:"default_merge_strategy,omitempty" mapstructure:"default_merge_strategy"`
}

type BranchReference struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

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

var columns = []string{
	"name",
	"target",
	"default_merge_strategy",
	"merge_strategies",
}

// NewReference creates a new BranchReference
func NewReference(name string) *BranchReference {
	return &BranchReference{
		Type: "branch",
		Name: name,
	}
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (branch Branch) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (branch Branch) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "name":
			row = append(row, branch.Name)
		case "target":
			row = append(row, branch.Target.Hash)
		case "default_merge_strategy", "default merge strategy":
			row = append(row, branch.DefaultMergeStrategy)
		case "merge_strategies", "merge strategies":
			row = append(row, strings.Join(branch.MergeStrategies, ", "))
		}
	}
	return row
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
