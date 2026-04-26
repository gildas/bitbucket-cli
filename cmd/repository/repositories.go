package repository

import (
	"context"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Repositories []Repository

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (repositories Repositories) GetHeaders(cmd *cobra.Command) []string {
	return Repository{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (repositories Repositories) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(repositories) {
		return []string{}
	}
	return repositories[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (repositories Repositories) Size() int {
	return len(repositories)
}

// GetRepositories gets the repositories for a workspace
func GetRepositories(ctx context.Context, cmd *cobra.Command) ([]Repository, error) {
	return GetRepositoriesWithQuery(ctx, cmd, url.Values{})
}

// GetRepositoriesWithQuery gets the repositories for a workspace with a query
func GetRepositoriesWithQuery(ctx context.Context, cmd *cobra.Command, query url.Values) ([]Repository, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("repository", "list")

	workspace, err := workspace.GetWorkspace(ctx, cmd)
	if err != nil {
		return nil, err
	}

	uriPath := "/repositories/" + workspace.Slug
	if len(query) > 0 {
		uriPath += "?" + query.Encode()
	}

	log.Infof("Getting repositories from workspace %s with query %s", workspace.Name, query)
	repositories, err := profile.GetAll[Repository](ctx, cmd, uriPath)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d repositories in workspace %s", len(repositories), workspace.Slug)
	core.Sort(repositories, func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Slug), strings.ToLower(b.Slug)) == -1
	})
	return repositories, nil
}

// GetRepositorySlugs gets the slugs of all repositories
func GetRepositorySlugs(ctx context.Context, cmd *cobra.Command) (slugs []string, err error) {
	repositories, err := GetRepositories(ctx, cmd)
	if err != nil {
		return
	}
	return core.Map(repositories, func(repository Repository) string { return repository.Slug }), nil
}

// GetRepositoryAllowedSlugs gets the slugs of all repositories to use with enum flags
func GetRepositoryAllowedSlugs(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) (slugs []string, err error) {
	return GetRepositorySlugs(ctx, cmd)
}
