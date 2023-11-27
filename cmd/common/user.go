package common

import "bitbucket.org/gildas_cherruel/bb/cmd/link"

type User struct {
	Type      string     `json:"type"          mapstructure:"type"`
	ID        string     `json:"uuid"          mapstructure:"uuid"`
	AccountID string     `json:"account_id"    mapstructure:"account_id"`
	Name      string     `json:"display_name"  mapstructure:"display_name"`
	Nickname  string     `json:"nickname"      mapstructure:"nickname"`
	Raw       string     `json:"raw,omitempty" mapstructure:"raw"`
	Links     link.Links `json:"links"         mapstructure:"links"`
}
