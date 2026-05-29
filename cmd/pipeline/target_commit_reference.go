package pipeline

import (
	"encoding/json"

	"github.com/gildas/bitbucket-cli/cmd/commit"
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-errors"
)

// CommitReferenceTarget represents the target of a pipeline (branch, tag, etc.)
type CommitReferenceTarget struct {
	Type     string                  `json:"type"               mapstructure:"type"`
	Selector *common.Selector        `json:"selector,omitempty" mapstructure:"selector"`
	Commit   *commit.CommitReference `json:"commit,omitempty"   mapstructure:"commit"`
}

func init() {
	targetRegistry.Add(CommitReferenceTarget{})
}

// GetType returns the target type
func (target CommitReferenceTarget) GetType() string {
	return "pipeline_commit_target"
}

// GetDestination returns the target's destination
//
// implements Target
func (target CommitReferenceTarget) GetDestination() string {
	return ""
}

// GetCommit return the target's commit reference
//
// implements Target
func (target CommitReferenceTarget) GetCommit() *commit.CommitReference {
	return target.Commit
}

// MarshalJSON custom JSON marshalling for Target
//
// implements json.Marshaler
func (target CommitReferenceTarget) MarshalJSON() ([]byte, error) {
	type surrogate CommitReferenceTarget

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
func (target *CommitReferenceTarget) UnmarshalJSON(data []byte) error {
	type surrogate CommitReferenceTarget
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
	*target = CommitReferenceTarget(inner.surrogate)

	return nil
}
