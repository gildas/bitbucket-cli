package pullrequest

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
)

// PullRequestReference describes a reference to a PullRequest
type PullRequestReference struct {
	ID      uint64       `json:"id"    mapstructure:"id"`
	Title   string       `json:"title,omitempty" mapstructure:"title"`
	IsDraft bool         `json:"draft,omitempty" mapstructure:"draft"`
	Links   common.Links `json:"links,omitempty" mapstructure:"links"`
}

// MarshalJSON marshals the PullRequestReference to JSON
//
// implements json.Marshaler
func (reference PullRequestReference) MarshalJSON() ([]byte, error) {
	type surrogate PullRequestReference
	var links *common.Links

	if !reference.Links.IsEmpty() {
		links = &reference.Links
	}

	data, err := json.Marshal(struct {
		surrogate
		Links *common.Links `json:"links,omitempty"`
	}{
		surrogate: surrogate(reference),
		Links:     links,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
