package pullrequest

import (
	"encoding/json"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-errors"
)

// PullRequestReference describes a reference to a PullRequest
type PullRequestReference struct {
	ID       uint64       `json:"id"               mapstructure:"id"`
	Title    string       `json:"title,omitempty"  mapstructure:"title"`
	IsDraft  bool         `json:"draft,omitempty"  mapstructure:"draft"`
	IsQueued bool         `json:"queued,omitempty" mapstructure:"queued"`
	Links    common.Links `json:"links,omitempty"  mapstructure:"links"`
}

// GetType returns the type of the PullRequestReference.
//
// implements core.TypeCarrier
func (reference PullRequestReference) GetType() string {
	return "pullrequest"
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
		Type string `json:"type"`
		surrogate
		Links *common.Links `json:"links,omitempty"`
	}{
		Type:      reference.GetType(),
		surrogate: surrogate(reference),
		Links:     links,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
