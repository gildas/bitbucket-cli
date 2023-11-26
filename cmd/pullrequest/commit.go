package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

type Commit struct {
	Type  string       `json:"type"  mapstructure:"type"`
	Hash  string       `json:"hash"  mapstructure:"hash"`
	Links common.Links `json:"links" mapstructure:"links"`
}
