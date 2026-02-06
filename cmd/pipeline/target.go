package pipeline

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
)

// validTargetTypes lists the target types accepted during unmarshal
var validTargetTypes = map[string]bool{
	"pipeline_ref_target":         true,
	"pipeline_pullrequest_target": true,
}

// Target represents the target of a pipeline (branch, tag, etc.)
type Target struct {
	Type              string           `json:"type"                          mapstructure:"type"`
	RefType           string           `json:"ref_type,omitempty"            mapstructure:"ref_type"`
	RefName           string           `json:"ref_name,omitempty"            mapstructure:"ref_name"`
	Selector          *Selector        `json:"selector,omitempty"            mapstructure:"selector"`
	Commit            *commit.Commit   `json:"commit,omitempty"              mapstructure:"commit"`
	Source            string           `json:"source,omitempty"              mapstructure:"source"`
	Destination       string           `json:"destination,omitempty"         mapstructure:"destination"`
	DestinationCommit *commit.Commit   `json:"destination_commit,omitempty"  mapstructure:"destination_commit"`
	PullRequest       *PullRequestRef  `json:"pullrequest,omitempty"         mapstructure:"pullrequest"`
}

// PullRequestRef represents a pull request reference in a pipeline target
type PullRequestRef struct {
	Type  string       `json:"type"            mapstructure:"type"`
	ID    int          `json:"id"              mapstructure:"id"`
	Title string       `json:"title,omitempty" mapstructure:"title"`
	Draft bool         `json:"draft,omitempty" mapstructure:"draft"`
	Links common.Links `json:"links,omitempty" mapstructure:"links"`
}

// Selector represents a pipeline selector for custom pipelines
type Selector struct {
	Type    string `json:"type"              mapstructure:"type"`
	Pattern string `json:"pattern,omitempty" mapstructure:"pattern"`
}

// GetType returns the target type
func (target Target) GetType() string {
	if target.Type != "" {
		return target.Type
	}
	return "pipeline_ref_target"
}

// MarshalJSON custom JSON marshalling for Target
//
// implements json.Marshaler
func (target Target) MarshalJSON() ([]byte, error) {
	type surrogate Target
	var commitRef *commit.CommitReference

	if target.Commit != nil {
		commitRef = target.Commit.GetReference()
	}

	var destCommitRef *commit.CommitReference
	if target.DestinationCommit != nil {
		destCommitRef = target.DestinationCommit.GetReference()
	}

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		CommitRef     *commit.CommitReference `json:"commit,omitempty"`
		DestCommitRef *commit.CommitReference `json:"destination_commit,omitempty"`
	}{
		Type:          target.GetType(),
		surrogate:     surrogate(target),
		CommitRef:     commitRef,
		DestCommitRef: destCommitRef,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON custom JSON unmarshalling for Target
//
// implements json.Unmarshaler
func (target *Target) UnmarshalJSON(data []byte) error {
	type surrogate Target
	var inner struct {
		Type string `json:"type"`
		surrogate
		CommitReference     *commit.CommitReference `json:"commit,omitempty"`
		DestCommitReference *commit.CommitReference `json:"destination_commit,omitempty"`
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if !validTargetTypes[inner.Type] {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, "pipeline_*_target"))
	}
	*target = Target(inner.surrogate)
	target.Type = inner.Type
	if inner.CommitReference != nil {
		target.Commit = inner.CommitReference.AsCommit()
	}
	if inner.DestCommitReference != nil {
		target.DestinationCommit = inner.DestCommitReference.AsCommit()
	}

	// For PR targets, populate RefName from source branch for display compatibility
	if inner.Type == "pipeline_pullrequest_target" && target.Source != "" {
		target.RefName = target.Source
	}

	return nil
}
