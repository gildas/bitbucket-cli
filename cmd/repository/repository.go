package repository

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
)

type Repository struct {
	Type     string      `json:"type"      mapstructure:"type"`
	ID       common.UUID `json:"uuid"      mapstructure:"uuid"`
	Name     string      `json:"name"      mapstructure:"name"`
	FullName string      `json:"full_name" mapstructure:"full_name"`
	Links    link.Links  `json:"links"     mapstructure:"links"`
}
