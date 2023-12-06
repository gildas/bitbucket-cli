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

type Member struct {
	Type      string     `json:"type"       mapstructure:"type"`
	User      user.User  `json:"user"       mapstructure:"user"`
	Workspace Workspace  `json:"workspace"  mapstructure:"workspace"`
	Links     link.Links `json:"links"      mapstructure:"links"`
}
