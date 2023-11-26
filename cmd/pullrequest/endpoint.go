package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
)

type Endpoint struct {
	Branch     Branch                `json:"branch"     mapstructure:"branch"`
	Commit     Commit                `json:"commit"     mapstructure:"commit"`
	Repository repository.Repository `json:"repository" mapstructure:"repository"`
}
