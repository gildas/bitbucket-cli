package workspace

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/spf13/cobra"
)

type Membership struct {
	Type       string       `json:"type"       mapstructure:"type"`
	Permission string       `json:"permission" mapstructure:"permission"`
	User       user.User    `json:"user"       mapstructure:"user"`
	Workspace  Workspace    `json:"workspace"  mapstructure:"workspace"`
	Links      common.Links `json:"links"      mapstructure:"links"`
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (membership Membership) GetHeaders(cmd *cobra.Command) []string {
	return []string{"ID", "Name", "Workspace", "Permission"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (membership Membership) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, membership.User.ID.String())
		case "name":
			row = append(row, membership.User.Name)
		case "workspace":
			row = append(row, membership.Workspace.Name)
		case "permission":
			row = append(row, membership.Permission)
		}
	}
	return row
}
