package workspace

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
)

type Membership struct {
	Type       string     `json:"type"       mapstructure:"type"`
	Permission string     `json:"permission" mapstructure:"permission"`
	User       user.User  `json:"user"       mapstructure:"user"`
	Workspace  Workspace  `json:"workspace"  mapstructure:"workspace"`
	Links      link.Links `json:"links"      mapstructure:"links"`
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (membership Membership) GetHeader(short bool) []string {
	return []string{"ID", "Name", "Workspace", "Permission"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (membership Membership) GetRow(headers []string) []string {
	return []string{
		membership.User.ID,
		membership.User.Name,
		membership.Workspace.Name,
		membership.Permission,
	}
}
