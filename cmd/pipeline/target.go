package pipeline

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"github.com/gildas/go-errors"
)

// Target represents the target of a pipeline (branch, tag, etc.)
type Target struct {
	Type     string         `json:"type"               mapstructure:"type"`
	RefType  string         `json:"ref_type,omitempty" mapstructure:"ref_type"`
	RefName  string         `json:"ref_name,omitempty" mapstructure:"ref_name"`
	Selector *Selector      `json:"selector,omitempty" mapstructure:"selector"`
	Commit   *commit.Commit `json:"commit,omitempty"   mapstructure:"commit"`
}

// Selector represents a pipeline selector for custom pipelines
type Selector struct {
	Type    string `json:"type"              mapstructure:"type"`
	Pattern string `json:"pattern,omitempty" mapstructure:"pattern"`
}

// GetType returns the target type
func (target Target) GetType() string {
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

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		CommitRef *commit.CommitReference `json:"commit,omitempty"`
	}{
		Type:      target.GetType(),
		surrogate: surrogate(target),
		CommitRef: commitRef,
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
		CommitReference *commit.CommitReference `json:"commit,omitempty"`
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != target.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, target.GetType()))
	}
	*target = Target(inner.surrogate)
	if inner.CommitReference != nil {
		target.Commit = inner.CommitReference.AsCommit()
	}

	return nil
}
