package pipeline

import (
	"encoding/json"

	"github.com/gildas/bitbucket-cli/cmd/commit"
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-errors"
)

// MergeQueueReferenceTarget represents the target of a pipeline run by a merge queue
type MergeQueueReferenceTarget struct {
	ReferenceType  string                  `json:"ref_type,omitempty"        mapstructure:"ref_type"`
	ReferenceName  string                  `json:"ref_name,omitempty"        mapstructure:"ref_name"`
	TargetBranch   string                  `json:"target_branch,omitempty"   mapstructure:"target_branch"`
	PullRequestIDs []uint64                `json:"pullrequest_ids,omitempty" mapstructure:"pullrequest_ids"`
	Selector       *common.Selector        `json:"selector,omitempty"        mapstructure:"selector"`
	Commit         *commit.CommitReference `json:"commit,omitempty"          mapstructure:"commit"`
}

func init() {
	targetRegistry.Add(MergeQueueReferenceTarget{})
}

// GetType returns the target type
func (target MergeQueueReferenceTarget) GetType() string {
	return "pipeline_merge_queue_ref_target"
}

// GetDestination returns the target's destination
//
// implements Target
func (target MergeQueueReferenceTarget) GetDestination() string {
	return target.TargetBranch
}

// GetCommit return the target's commit reference
//
// implements Target
func (target MergeQueueReferenceTarget) GetCommit() *commit.CommitReference {
	return target.Commit
}

// MarshalJSON custom JSON marshalling for Target
//
// implements json.Marshaler
func (target MergeQueueReferenceTarget) MarshalJSON() ([]byte, error) {
	type surrogate MergeQueueReferenceTarget

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
func (target *MergeQueueReferenceTarget) UnmarshalJSON(data []byte) error {
	type surrogate MergeQueueReferenceTarget
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
	*target = MergeQueueReferenceTarget(inner.surrogate)

	return nil
}
