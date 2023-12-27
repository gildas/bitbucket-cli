package repository

import (
	"context"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Repository struct {
	Type     string       `json:"type"      mapstructure:"type"`
	ID       common.UUID  `json:"uuid"      mapstructure:"uuid"`
	Name     string       `json:"name"      mapstructure:"name"`
	FullName string       `json:"full_name" mapstructure:"full_name"`
	Links    common.Links `json:"links"     mapstructure:"links"`
}

func GetRepository(context context.Context, cmd *cobra.Command) (*Repository, error) {
	fullName := cmd.Flag("repository").Value.String()
	if len(fullName) == 0 {
		remote, err := remote.GetFromGitConfig(context, "origin")
		if err != nil {
			return nil, errors.Join(errors.NotFound.With("current repository"), err)
		}
		fullName = remote.RepositoryName()
	}
	components := strings.Split(fullName, "/")
	if len(components) == 2 {
		return &Repository{Name: components[1], FullName: fullName}, nil
	} else if len(components) == 1 {
		return &Repository{Name: components[0], FullName: fullName}, nil
	}
	return &Repository{Name: fullName, FullName: fullName}, nil
}

// String returns the string representation of the repository
//
// implements fmt.Stringer
func (repository Repository) String() string {
	return repository.FullName
}
