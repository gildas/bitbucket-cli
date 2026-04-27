package commit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type Commit struct {
	Hash       string                `json:"hash"               mapstructure:"hash"`
	Author     user.Author           `json:"author"             mapstructure:"author"`
	Message    string                `json:"message"            mapstructure:"message"`
	Summary    *common.RenderedText  `json:"summary,omitempty"  mapstructure:"summary"`
	Rendered   *RenderedMessage      `json:"rendered,omitempty" mapstructure:"rendered"`
	Parents    []CommitReference     `json:"parents,omitempty"  mapstructure:"parents"`
	Date       time.Time             `json:"date"               mapstructure:"date"`
	Repository repository.Repository `json:"repository"         mapstructure:"repository"`
	Links      common.Links          `json:"links"              mapstructure:"links"`
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
	{Name: "hash", DefaultSorter: false, Compare: func(a, b Commit) bool {
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
	{Name: "date", DefaultSorter: true, Compare: func(a, b Commit) bool {
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
		Hash:  commit.Hash,
		Links: commit.Links,
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
	return []string{"Hash", "Date", "Author", "Message"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (commit Commit) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "hash":
			row = append(row, commit.GetShortHash())
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

// GetShortHash gets the short hash of this commit
func (commit Commit) GetShortHash() string {
	if len(commit.Hash) > 7 {
		return commit.Hash[:7]
	}
	return commit.Hash
}

// GetLatestCommit gets the latest commit of the repository
func GetLatestCommit(ctx context.Context, cmd *cobra.Command) (commit *Commit, err error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return GetCommitByHash(ctx, cmd, head.Hash().String())
}

// GetCommitByHash gets a commit by its hash
func GetCommitByHash(ctx context.Context, cmd *cobra.Command, hash string) (commit *Commit, err error) {
	profile, err := profile.GetProfileFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}
	err = profile.Get(ctx, cmd, repository.GetPath("commit", hash), &commit)
	return commit, err
}

// Validate validates a Commit
func (commit *Commit) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this commit
//
// implements fmt.Stringer
func (commit Commit) String() string {
	return commit.Hash
}

// MarshalJSON implements the json.Marshaler interface.
func (commit Commit) MarshalJSON() (data []byte, err error) {
	type surrogate Commit

	data, err = json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		Date string `json:"date"`
	}{
		Type:      commit.GetType(),
		surrogate: surrogate(commit),
		Date:      commit.Date.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
