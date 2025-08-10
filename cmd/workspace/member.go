package workspace

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/spf13/cobra"
)

type Member struct {
	Type      string       `json:"type"      mapstructure:"type"`
	User      user.User    `json:"user"      mapstructure:"user"`
	Workspace Workspace    `json:"workspace" mapstructure:"workspace"`
	Links     common.Links `json:"links"     mapstructure:"links"`
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (member Member) GetHeaders(cmd *cobra.Command) []string {
	return []string{"ID", "Name", "Workspace"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (member Member) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, member.User.ID.String())
		case "name":
			row = append(row, member.User.Name)
		case "workspace":
			row = append(row, member.Workspace.Name)
		}
	}
	return row
}
