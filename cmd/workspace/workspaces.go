package workspace

import (
	"context"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Workspaces []Workspace

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Slug"}
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (workspaces Workspaces) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(workspaces) {
		return []string{}
	}
	return workspaces[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (workspaces Workspaces) Size() int {
	return len(workspaces)
}

// GetWorkspaces gets the workspaces for the current user
func GetWorkspaces(ctx context.Context, cmd *cobra.Command) (Workspaces, error) {
	return GetWorkspacesWithQuery(ctx, cmd, url.Values{})
}

// GetWorkspacesWithQuery gets the workspaces for the current user with a query
func GetWorkspacesWithQuery(ctx context.Context, cmd *cobra.Command, query url.Values) (Workspaces, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("workspace", "slugs")

	uripath := "/user/workspaces"
	if len(query) > 0 {
		uripath += "?" + query.Encode()
	}

	log.Debugf("Getting all workspaces with query %s", query)
	workspaces, err := profile.GetAll[Workspace](ctx, cmd, uripath)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d workspaces", len(workspaces))
	core.Sort(workspaces, func(a, b Workspace) bool {
		return strings.Compare(strings.ToLower(a.Slug), strings.ToLower(b.Slug)) == -1
	})
	return workspaces, nil
}

// GetWorkspaceSlugs gets the slugs of all workspaces
func GetWorkspaceSlugs(ctx context.Context, cmd *cobra.Command) (slugs []string, err error) {
	workspaces, err := GetWorkspaces(ctx, cmd)
	if err != nil {
		return
	}
	return core.Map(workspaces, func(workspace Workspace) string { return workspace.Slug }), nil
}

// GetWorkspaceAllowedSlugs gets the slugs of all workspaces to use with enum flag completion
func GetWorkspaceAllowedSlugs(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) (slugs []string, err error) {
	return GetWorkspaceSlugs(ctx, cmd)
}
