package commit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Commit struct {
	Type       string                `json:"type"               mapstructure:"type"`
	Hash       string                `json:"hash"               mapstructure:"hash"`
	Author     user.Author           `json:"author"             mapstructure:"author"`
	Message    string                `json:"message"            mapstructure:"message"`
	Summary    *common.RenderedText  `json:"summary,omitempty"  mapstructure:"summary"`
	Rendered   *RenderedMessage      `json:"rendered,omitempty" mapstructure:"rendered"`
	Parents    []CommitRef           `json:"parents"            mapstructure:"parents"`
	Date       time.Time             `json:"date"               mapstructure:"date"`
	Repository repository.Repository `json:"repository"         mapstructure:"repository"`
	Links      common.Links          `json:"links"              mapstructure:"links"`
}

type CommitRef struct {
	Type  string       `json:"type"  mapstructure:"type"`
	Hash  string       `json:"hash"  mapstructure:"hash"`
	Links common.Links `json:"links" mapstructure:"links"`
}

type RenderedMessage struct {
	Message common.RenderedText `json:"message" mapstructure:"message"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "commit",
	Short: "Manage commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Commit requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Commit]{
	{Name: "hash", DefaultSorter: true, Compare: func(a, b Commit) bool {
		return strings.Compare(strings.ToLower(a.Hash), strings.ToLower(b.Hash)) == -1
	}},
	{Name: "longhash", DefaultSorter: false, Compare: func(a, b Commit) bool {
		return strings.Compare(strings.ToLower(a.Message), strings.ToLower(b.Message)) == -1
	}},
	{Name: "author", DefaultSorter: false, Compare: func(a, b Commit) bool {
		return strings.Compare(strings.ToLower(a.Author.User.Name), strings.ToLower(b.Author.User.Name)) == -1
	}},
	{Name: "message", DefaultSorter: false, Compare: func(a, b Commit) bool {
		return strings.Compare(strings.ToLower(a.Message), strings.ToLower(b.Message)) == -1
	}},
	{Name: "date", DefaultSorter: false, Compare: func(a, b Commit) bool {
		return a.Date.Before(b.Date)
	}},
	{Name: "repository", DefaultSorter: false, Compare: func(a, b Commit) bool {
		return strings.Compare(strings.ToLower(a.Repository.Name), strings.ToLower(b.Repository.Name)) == -1
	}},
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (commit Commit) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Hash", "Author", "Message"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (commit Commit) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "hash":
			row = append(row, commit.Hash[:7])
		case "longhash", "fullhash":
			row = append(row, commit.Hash)
		case "author":
			row = append(row, commit.Author.User.Name)
		case "message":
			row = append(row, commit.Message)
		case "date":
			row = append(row, commit.Date.Format("2006-01-02 15:04:05"))
		case "repository":
			row = append(row, commit.Repository.Name)
		}
	}
	return row
}

// Validate validates a Commit
func (commit *Commit) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (commit Commit) String() string {
	return commit.Hash
}

// MarshalJSON implements the json.Marshaler interface.
func (commit Commit) MarshalJSON() (data []byte, err error) {
	type surrogate Commit

	data, err = json.Marshal(struct {
		surrogate
		Date string `json:"date"`
	}{
		surrogate: surrogate(commit),
		Date:      commit.Date.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
