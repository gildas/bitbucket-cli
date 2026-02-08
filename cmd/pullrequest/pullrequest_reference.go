package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

// PullRequestReference describe a reference to a PullRequest
type PullRequestReference struct {
	Type    string       `json:"type"  mapstructure:"type"`
	ID      uint64       `json:"id"    mapstructure:"id"`
	Title   string       `json:"title" mapstructure:"title"`
	IsDraft bool         `json:"draft" mapstructure:"draft"`
	Links   common.Links `json:"links" mapstructure:"links"`
}
