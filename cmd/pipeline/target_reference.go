package pipeline

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
)

// ReferenceTarget represents the target of a pipeline (branch, tag, etc.)
type ReferenceTarget struct {
	Type          string                 `json:"type"               mapstructure:"type"`
	ReferenceType string                 `json:"ref_type,omitempty" mapstructure:"ref_type"`
	ReferenceName string                 `json:"ref_name,omitempty" mapstructure:"ref_name"`
	Selector      *common.Selector       `json:"selector,omitempty" mapstructure:"selector"`
	Commit        commit.CommitReference `json:"commit"             mapstructure:"commit"`
}

func init() {
	targetRegistry.Add(ReferenceTarget{})
}

// GetType returns the target type
func (target ReferenceTarget) GetType() string {
	return "pipeline_ref_target"
}

// GetDestination returns the target's destination
//
// implements Target
func (target ReferenceTarget) GetDestination() string {
	return target.ReferenceName
}

// GetCommit return the target's commit reference
func (target ReferenceTarget) GetCommit() commit.CommitReference {
	return target.Commit
}

// MarshalJSON custom JSON marshalling for Target
//
// implements json.Marshaler
func (target ReferenceTarget) MarshalJSON() ([]byte, error) {
	type surrogate ReferenceTarget

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      target.GetType(),
		surrogate: surrogate(target),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON custom JSON unmarshalling for Target
//
// implements json.Unmarshaler
func (target *ReferenceTarget) UnmarshalJSON(data []byte) error {
	type surrogate ReferenceTarget
	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != target.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, target.GetType()))
	}
	*target = ReferenceTarget(inner.surrogate)

	return nil
}
