package branch

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type Branch struct {
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

var columns = common.Columns[Branch]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b Branch) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "target", DefaultSorter: false, Compare: func(a, b Branch) bool {
		return strings.Compare(strings.ToLower(a.Target.Hash), strings.ToLower(b.Target.Hash)) == -1
	}},
	{Name: "default_merge_strategy", DefaultSorter: false, Compare: func(a, b Branch) bool {
		return strings.Compare(strings.ToLower(a.DefaultMergeStrategy), strings.ToLower(b.DefaultMergeStrategy)) == -1
	}},
	{Name: "merge_strategies", DefaultSorter: false, Compare: func(a, b Branch) bool {
		return strings.Compare(strings.ToLower(strings.Join(a.MergeStrategies, ",")), strings.ToLower(strings.Join(b.MergeStrategies, ","))) == -1
	}},
}

// GetCurrentBranch gets the current branch from Git repository
func GetCurrentBranch() (*Branch, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	if !head.Name().IsBranch() {
		return nil, errors.Errorf("current HEAD is not a branch (%s)", head.Name().String())
	}
	return &Branch{Name: head.Name().Short()}, nil
}

// GetType returns the branch type
func (branch Branch) GetType() string {
	return "branch"
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

// MarshalJSON custom JSON marshalling for Branch
//
// implements json.Marshaler
func (branch Branch) MarshalJSON() ([]byte, error) {
	type surrogate Branch
	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      branch.GetType(),
		surrogate: surrogate(branch),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON custom JSON unmarshalling for Branch
//
// implements json.Unmarshaler
func (branch *Branch) UnmarshalJSON(data []byte) error {
	type surrogate Branch
	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != branch.GetType() {
		return errors.JSONUnmarshalError.With("invalid type: expected %s, got %s", branch.GetType(), inner.Type)
	}
	*branch = Branch(inner.surrogate)
	return nil
}
