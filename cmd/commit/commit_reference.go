package commit

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
)

type CommitReference struct {
	Hash  string       `json:"hash"  mapstructure:"hash"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// AsCommit converts this CommitRef to a Commit
func (reference CommitReference) AsCommit() *Commit {
	return &Commit{
		Hash:  reference.Hash,
		Links: reference.Links,
	}
}

// GetShortHash gets the short hash of this commit
func (reference CommitReference) GetShortHash() string {
	if len(reference.Hash) > 7 {
		return reference.Hash[:7]
	}
	return reference.Hash
}

// String gets a string representation of this commit
//
// implements fmt.Stringer
func (reference CommitReference) String() string {
	return reference.Hash
}

// MarshalJSON implements the json.Marshaler interface.
func (reference CommitReference) MarshalJSON() (data []byte, err error) {
	type surrogate CommitReference
	var links *common.Links

	if !reference.Links.IsEmpty() {
		links = &reference.Links
	}

	data, err = json.Marshal(struct {
		// Type string `json:"type"`
		surrogate
		Links *common.Links `json:"links,omitempty"`
	}{
		// Type:      "commit",
		surrogate: surrogate(reference),
		Links:     links,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
