package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Workspace struct {
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

var columns = common.Columns[Workspace]{
	{Name: "id", DefaultSorter: false, Compare: func(a, b Workspace) bool {
		return strings.Compare(strings.ToLower(a.ID.String()), strings.ToLower(b.ID.String())) == -1
	}},
	{Name: "name", DefaultSorter: true, Compare: func(a, b Workspace) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "slug", DefaultSorter: false, Compare: func(a, b Workspace) bool {
		return strings.Compare(strings.ToLower(a.Slug), strings.ToLower(b.Slug)) == -1
	}},
}

var WorkspaceCache = common.NewCache[Workspace]()

// GetType gets the type of the workspace
//
// implements core.TypeCarrier
func (workspace Workspace) GetType() string {
	return "workspace"
}

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

// String returns the string representation of the workspace
//
// implements fmt.Stringer
func (workspace Workspace) String() string {
	return workspace.Slug
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

// MarshalJSON marshals the workspace to JSON
//
// implements json.Marshaler
func (workspace Workspace) MarshalJSON() ([]byte, error) {
	type surrogate Workspace

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      workspace.GetType(),
		surrogate: surrogate(workspace),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the workspace from JSON
//
// implements json.Unmarshaler
func (workspace *Workspace) UnmarshalJSON(data []byte) error {
	type surrogate Workspace

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != workspace.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, workspace.GetType()))
	}

	*workspace = Workspace(inner.surrogate)
	return nil
}
