package pipeline

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest"
)

// PullRequestReferenceTarget represents a target for a pipeline that is a pull request reference.
type PullRequestReferenceTarget struct {
	Source            string                           `json:"source" mapstructure:"source"`
	Destination       string                           `json:"destination" mapstructure:"destination"`
	DestinationCommit commit.CommitReference           `json:"destination_commit" mapstructure:"destination_commit"`
	Commit            commit.CommitReference           `json:"commit" mapstructure:"commit"`
	Selector          *common.Selector                 `json:"selector" mapstructure:"selector"`
	PullRequest       pullrequest.PullRequestReference `json:"pullrequest" mapstructure:"pullrequest"`
}

func init() {
	targetRegistry.Add(PullRequestReferenceTarget{})
}

// GetType returns the type of the PullRequestReferenceTarget.
//
// implements core.TypeCarier
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
func (target PullRequestReferenceTarget) GetCommit() commit.CommitReference {
	return target.Commit
}
