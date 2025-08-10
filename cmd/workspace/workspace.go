package workspace

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Workspace struct {
	Type  string       `json:"type"  mapstructure:"type"`
	ID    common.UUID  `json:"uuid"  mapstructure:"uuid"`
	Name  string       `json:"name"  mapstructure:"name"`
	Slug  string       `json:"slug"  mapstructure:"slug"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = []string{
	"id",
	"name",
	"slug",
}

var WorkspaceCache = common.NewCache[Workspace]()

// GetID gets the ID of the workspace
//
// implements core.Identifiable
func (workspace Workspace) GetID() uuid.UUID {
	return uuid.UUID(workspace.ID)
}

// GetName gets the name of the workspace
//
// implements core.Named
func (workspace Workspace) GetName() string {
	return workspace.Name
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (workspace Workspace) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Name", "Slug"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (workspace Workspace) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, workspace.ID.String())
		case "name":
			row = append(row, workspace.Name)
		case "slug":
			row = append(row, workspace.Slug)
		}
	}
	return row
}

// GetWorkspace gets the workspace by its slug
func GetWorkspace(context context.Context, cmd *cobra.Command, workspaceName string) (workspace *Workspace, err error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "get")

	currentProfile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}

	log.Infof("Retrieving workspace %s", workspaceName)

	if workspace, err = WorkspaceCache.Get(workspaceName); err == nil {
		log.Debugf("Workspace %s found in cache", workspaceName)
		return workspace, nil
	}

	err = currentProfile.Get(
		log.ToContext(context),
		cmd,
		fmt.Sprintf("/workspaces/%s", workspaceName),
		&workspace,
	)
	if err == nil {
		_ = WorkspaceCache.Set(*workspace, workspaceName)
	}

	return
}

// GetWorkspaceFromGit gets the workspace from the git config
func GetWorkspaceFromGit(context context.Context, cmd *cobra.Command) (workspace *Workspace, err error) {
	remote, err := remote.GetFromGitConfig(context, "origin")
	if err != nil {
		return nil, err
	}
	return GetWorkspace(context, cmd, remote.WorkspaceName())
}

// GetMembers gets the members of the workspace
func (workspace Workspace) GetMembers(context context.Context, cmd *cobra.Command) (members []Member, err error) {
	members, err = profile.GetAll[Member](
		cmd.Context(),
		cmd,
		fmt.Sprintf("/workspaces/%s/members", workspace.Slug),
	)
	if err != nil {
		return []Member{}, err
	}
	return
}

// GetWorkspaceSlugs gets the slugs of all workspaces
func GetWorkspaceSlugs(context context.Context, cmd *cobra.Command, args []string, toComplete string) (slugs []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "slugs")

	log.Debugf("Getting all workspaces")
	workspaces, err := profile.GetAll[Workspace](context, cmd, "/workspaces")
	if err != nil {
		log.Errorf("Failed to get workspaces", err)
		return
	}
	return core.Map(workspaces, func(workspace Workspace) string {
		return workspace.Slug
	}), nil
}
