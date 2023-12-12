package workspace

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
)

type Member struct {
	Type      string     `json:"type"       mapstructure:"type"`
	User      user.User  `json:"user"       mapstructure:"user"`
	Workspace Workspace  `json:"workspace"  mapstructure:"workspace"`
	Links     link.Links `json:"links"      mapstructure:"links"`
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (member Member) GetHeader(short bool) []string {
	return []string{"ID", "Name", "Workspace"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (member Member) GetRow(headers []string) []string {
	return []string{
		member.User.ID,
		member.User.Name,
		member.Workspace.Name,
	}
}
