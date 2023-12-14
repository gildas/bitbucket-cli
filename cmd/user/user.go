package user

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
)

type User struct {
	Type      string      `json:"type"          mapstructure:"type"`
	ID        common.UUID `json:"uuid"          mapstructure:"uuid"`
	AccountID string      `json:"account_id"    mapstructure:"account_id"`
	Name      string      `json:"display_name"  mapstructure:"display_name"`
	Nickname  string      `json:"nickname"      mapstructure:"nickname"`
	Raw       string      `json:"raw,omitempty" mapstructure:"raw"`
	Links     link.Links  `json:"links"         mapstructure:"links"`
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (user User) GetHeader(short bool) []string {
	return []string{"ID", "Name", "Nickname"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (user User) GetRow(headers []string) []string {
	return []string{user.ID.String(), user.Name, user.Nickname}
}
