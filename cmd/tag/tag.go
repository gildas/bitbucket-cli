package tag

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

// Tag represents a Bitbucket tag
type Tag struct {
	Name    string        `json:"name"              mapstructure:"name"`
	Message string        `json:"message,omitempty" mapstructure:"message"`
	Author  user.Author   `json:"tagger"            mapstructure:"tagger"`
	Commit  commit.Commit `json:"target"            mapstructure:"target"`
	Date    time.Time     `json:"date"              mapstructure:"date"`
	Links   common.Links  `json:"links"             mapstructure:"links"`
}

type TagReference struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tag requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Tag]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b Tag) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "message", DefaultSorter: false, Compare: func(a, b Tag) bool {
		return strings.Compare(strings.ToLower(a.Message), strings.ToLower(b.Message)) == -1
	}},
	{Name: "author", DefaultSorter: false, Compare: func(a, b Tag) bool {
		return strings.Compare(strings.ToLower(a.Author.User.Name), strings.ToLower(b.Author.User.Name)) == -1
	}},
	{Name: "commit", DefaultSorter: false, Compare: func(a, b Tag) bool {
		return strings.Compare(strings.ToLower(a.Commit.Hash), strings.ToLower(b.Commit.Hash)) == -1
	}},
	{Name: "longcommit", DefaultSorter: false, Compare: func(a, b Tag) bool {
		return strings.Compare(strings.ToLower(a.Commit.Hash), strings.ToLower(b.Commit.Hash)) == -1
	}},
	{Name: "date", DefaultSorter: false, Compare: func(a, b Tag) bool {
		return a.Date.Before(b.Date)
	}},
}

// GetType returns the tag type
func (tag Tag) GetType() string {
	return "tag"
}

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (tag Tag) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"name", "commit", "date", "author", "message"}
}

// GetRow gets the row for a table
//
// implements common.Tableables
func (tag Tag) GetRow(headers []string) []string {
	row := make([]string, len(headers))
	for i, header := range headers {
		switch header {
		case "name":
			row[i] = tag.Name
		case "author":
			row[i] = tag.Author.User.Name
		case "commit":
			row[i] = tag.Commit.GetShortHash()
		case "longcommit":
			row[i] = tag.Commit.Hash
		case "date":
			row[i] = tag.Date.Format(time.RFC3339)
		case "message":
			row[i] = strings.TrimSpace(tag.Message)
		default:
			row[i] = ""
		}
	}
	return row
}

// UnmashalJSON unmarshals the tag from JSON
//
// implements json.Unmarshaler
func (tag *Tag) UnmarshalJSON(data []byte) error {
	type surrogate Tag
	var inner struct {
		Type string    `json:"type"`
		Date core.Time `json:"date"`
		surrogate
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != tag.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, tag.GetType()))
	}
	*tag = Tag(inner.surrogate)
	tag.Date = time.Time(inner.Date)
	return nil
}
