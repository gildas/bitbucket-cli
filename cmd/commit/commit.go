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
	Hash       string                `json:"hash"               mapstructure:"hash"`
	Author     user.Author           `json:"author"             mapstructure:"author"`
	Message    string                `json:"message"            mapstructure:"message"`
	Summary    *common.RenderedText  `json:"summary,omitempty"  mapstructure:"summary"`
	Rendered   *RenderedMessage      `json:"rendered,omitempty" mapstructure:"rendered"`
	Parents    []CommitReference     `json:"parents"            mapstructure:"parents"`
	Date       time.Time             `json:"date"               mapstructure:"date"`
	Repository repository.Repository `json:"repository"         mapstructure:"repository"`
	Links      common.Links          `json:"links"              mapstructure:"links"`
}

type CommitReference struct {
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

// GetType gets the type of this commit
//
// implements core.TypeCarrier
func (commit Commit) GetType() string {
	return "commit"
}

// GetReference gets the reference string for this commit
func (commit Commit) GetReference() *CommitReference {
	return &CommitReference{
		Type:  commit.GetType(),
		Hash:  commit.Hash,
		Links: commit.Links,
	}
}

// AsCommit converts this CommitRef to a Commit
func (reference CommitReference) AsCommit() *Commit {
	return &Commit{
		Hash:  reference.Hash,
		Links: reference.Links,
	}
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
	var author *user.Author
	var repo *repository.Repository
	var parents []CommitReference
	var date string

	if commit.Author.Type != "" || commit.Author.Raw != "" {
		author = &commit.Author
	}
	if commit.Repository.Name != "" || !commit.Repository.ID.IsNil() {
		repo = &commit.Repository
	}
	if len(commit.Parents) > 0 {
		parents = commit.Parents
	}
	if !commit.Date.IsZero() {
		date = commit.Date.Format("2006-01-02T15:04:05.999999999-07:00")
	}

	data, err = json.Marshal(struct {
		Type       string                `json:"type"`
		Hash       string                `json:"hash"`
		Author     *user.Author          `json:"author,omitempty"`
		Message    string                `json:"message,omitempty"`
		Summary    *common.RenderedText  `json:"summary,omitempty"`
		Rendered   *RenderedMessage      `json:"rendered,omitempty"`
		Parents    []CommitReference     `json:"parents,omitempty"`
		Date       string                `json:"date,omitempty"`
		Repository *repository.Repository `json:"repository,omitempty"`
		Links      common.Links          `json:"links"`
	}{
		Type:       commit.GetType(),
		Hash:       commit.Hash,
		Author:     author,
		Message:    commit.Message,
		Summary:    commit.Summary,
		Rendered:   commit.Rendered,
		Parents:    parents,
		Date:       date,
		Repository: repo,
		Links:      commit.Links,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// MarshalJSON implements the json.Marshaler interface.
func (ref CommitReference) MarshalJSON() (data []byte, err error) {
	type surrogate CommitReference
	var links *common.Links

	if !ref.Links.IsEmpty() {
		links = &ref.Links
	}

	data, err = json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		Links *common.Links `json:"links,omitempty"`
	}{
		Type:      "commit",
		surrogate: surrogate(ref),
		Links:     links,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
