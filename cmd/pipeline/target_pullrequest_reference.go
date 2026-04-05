package pipeline

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest"
	"github.com/gildas/go-errors"
)

// PullRequestReferenceTarget represents a target for a pipeline that is a pull request reference.
type PullRequestReferenceTarget struct {
	Source            string                           `json:"source"                       mapstructure:"source"`
	Destination       string                           `json:"destination"                  mapstructure:"destination"`
	DestinationCommit *commit.CommitReference          `json:"destination_commit,omitempty" mapstructure:"destination_commit"`
	Commit            *commit.CommitReference          `json:"commit,omitempty"             mapstructure:"commit"`
	Selector          *common.Selector                 `json:"selector,omitempty"           mapstructure:"selector"`
	PullRequest       pullrequest.PullRequestReference `json:"pullrequest"                  mapstructure:"pullrequest"`
}

func init() {
	targetRegistry.Add(PullRequestReferenceTarget{})
}

// GetType returns the type of the PullRequestReferenceTarget.
//
// implements core.TypeCarrier
func (target PullRequestReferenceTarget) GetType() string {
	return "pipeline_pullrequest_target"
}

// GetDestination returns the destination of the PullRequestReferenceTarget.
//
// implements Target
func (target PullRequestReferenceTarget) GetDestination() string {
	return target.Destination
}

// GetCommit returns the commit of the PullRequestReferenceTarget.
//
// implements Target
func (target PullRequestReferenceTarget) GetCommit() *commit.CommitReference {
	return target.Commit
}

// MarshalJSON marshals the PullRequestReferenceTarget to JSON
//
// implements json.Marshaler
func (target PullRequestReferenceTarget) MarshalJSON() ([]byte, error) {
	type surrogate PullRequestReferenceTarget

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      target.GetType(),
		surrogate: surrogate(target),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
