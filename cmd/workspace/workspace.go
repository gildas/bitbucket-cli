package workspace

import (
	"context"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
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

// GetHeader gets the header for a table
//
// implements common.Tableable
func (workspace Workspace) GetHeader(short bool) []string {
	return []string{"ID", "Name", "Slug"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (workspace Workspace) GetRow(headers []string) []string {
	return []string{
		workspace.ID.String(),
		workspace.Name,
		workspace.Slug,
	}
}

// GetWorkspace gets the workspace by its slug
func GetWorkspace(context context.Context, cmd *cobra.Command, profile *profile.Profile, workspace string) (*Workspace, error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "get")

	if profile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	log.Infof("Retrieving workspace %s", workspace)
	var result Workspace

	err := profile.Get(
		log.ToContext(context),
		cmd,
		fmt.Sprintf("/workspaces/%s", workspace),
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWorkspaceFromGit gets the workspace from the git config
func GetWorkspaceFromGit(context context.Context, cmd *cobra.Command, profile *profile.Profile) (workspace *Workspace, err error) {
	remote, err := remote.GetFromGitConfig(context, "origin")
	if err != nil {
		return nil, err
	}
	return GetWorkspace(context, cmd, profile, remote.WorkspaceName())
}

// GetMembers gets the members of the workspace
func (workspace Workspace) GetMembers(context context.Context, cmd *cobra.Command) (members []Member, err error) {
	members, err = profile.GetAll[Member](
		cmd.Context(),
		cmd,
		profile.Current,
		fmt.Sprintf("/workspaces/%s/members", workspace.Slug),
	)
	if err != nil {
		return []Member{}, err
	}
	return
}

// GetWorkspaceSlugs gets the slugs of all workspaces
func GetWorkspaceSlugs(context context.Context, cmd *cobra.Command, args []string) (slugs []string) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "slugs")

	log.Debugf("Getting all workspaces")
	workspaces, err := profile.GetAll[Workspace](context, cmd, profile.Current, "/workspaces")
	if err != nil {
		log.Errorf("Failed to get workspaces", err)
		return
	}
	return core.Map(workspaces, func(workspace Workspace) string {
		return workspace.Slug
	})
}
