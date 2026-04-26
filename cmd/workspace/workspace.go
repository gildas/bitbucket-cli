package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace/permission"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Workspace struct {
	ID            common.UUID  `json:"uuid"  mapstructure:"uuid"`
	Name          string       `json:"name"  mapstructure:"name"`
	Slug          string       `json:"slug"  mapstructure:"slug"`
	Administrator bool         `json:"administrator" mapstructure:"administrator"`
	Links         common.Links `json:"links" mapstructure:"links"`
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

func init() {
	Command.AddCommand(permission.Command)
}

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
	if len(workspace.Slug) > 0 {
		return workspace.Slug
	}
	return workspace.Name
}

// GetWorkspaceName gets the workspace name from the command flag or git config
//
// The workspace is determined by the following order:
//  1. The workspace flag in the command
//  2. The git config
//  3. The default workspace in the profile
func GetWorkspaceName(context context.Context, cmd *cobra.Command) (workspaceName string, err error) {
	if cmd.Flag("workspace") != nil {
		if workspaceName = cmd.Flag("workspace").Value.String(); len(workspaceName) > 0 {
			return
		}
	}
	if remote, err := remote.GetRemote(context, cmd); err == nil {
		return remote.WorkspaceName(), nil
	}
	if profile.Current != nil && len(profile.Current.DefaultWorkspace) > 0 {
		return profile.Current.DefaultWorkspace, nil
	}
	return "", errors.ArgumentMissing.With("workspace")
}

// GetWorkspace gets the current workspace
//
// The workspace is determined by the following order:
// 1. The workspace flag in the command
// 2. The git config
// 3. The default workspace in the profile
func GetWorkspace(ctx context.Context, cmd *cobra.Command) (workspace *Workspace, err error) {
	workspaceName, err := GetWorkspaceName(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return GetWorkspaceBySlugOrID(ctx, cmd, workspaceName)
}

// GetWorkspaceBySlugOrID gets the workspace by its slug name or ID
func GetWorkspaceBySlugOrID(ctx context.Context, cmd *cobra.Command, slugOrID string) (workspace *Workspace, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("workspace", "get_by_slug_or_id", "workspace", slugOrID)

	currentProfile, err := profile.GetProfileFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	log.Infof("Retrieving workspace %s", slugOrID)

	// In case we got a real UUID, get the Bitbucket UUID
	if id, err := common.ParseUUID(slugOrID); err == nil {
		slugOrID = id.String()
	}

	if workspace, err = WorkspaceCache.Get(slugOrID); err == nil {
		log.Debugf("Workspace %s found in cache", slugOrID)
		return workspace, nil
	}

	err = currentProfile.Get(
		log.ToContext(ctx),
		cmd,
		fmt.Sprintf("/workspaces/%s", slugOrID),
		&workspace,
	)
	if err == nil {
		_ = WorkspaceCache.Set(*workspace, slugOrID)
	}
	return workspace, errors.Join(errors.Errorf("Failed to get workspace %s", slugOrID), err)
}

// GetMember gets the workspace member by its username
func (workspace Workspace) GetMember(ctx context.Context, cmd *cobra.Command, profile *profile.Profile, username string) (member *Member, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("workspace", "get", "member", username)

	log.Infof("Retrieving workspace %s member %s", workspace.Slug, username)
	err = profile.Get(
		log.ToContext(ctx),
		cmd,
		fmt.Sprintf("/workspaces/%s/members/%s", workspace.Slug, username),
		&member,
	)
	return
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

	var typeholder struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeholder); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	switch typeholder.Type {
	case "workspace_access":
		var inner struct {
			Type          string `json:"type"`
			Administrator bool   `json:"administrator"`
			Workspace     struct {
				Type string `json:"type"`
				surrogate
			} `json:"workspace"`
		}
		if err := json.Unmarshal(data, &inner); err != nil {
			return errors.JSONUnmarshalError.WrapIfNotMe(err)
		}
		if inner.Workspace.Type != "workspace_base" {
			return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Workspace.Type, "workspace_base"))
		}

		*workspace = Workspace(inner.Workspace.surrogate)
		workspace.Administrator = inner.Administrator
	case "workspace":
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
	default:
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(typeholder.Type, strings.Join([]string{Workspace{}.GetType(), "workspace_access"}, ", ")))
	}

	return nil
}
