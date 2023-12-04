package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
)

type Endpoint struct {
	Branch     Branch                 `json:"branch"               mapstructure:"branch"`
	Commit     *commit.Commit         `json:"commit,omitempty"     mapstructure:"commit"`
	Repository *repository.Repository `json:"repository,omitempty" mapstructure:"repository"`
}
