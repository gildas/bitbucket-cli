package project

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project/reviewer"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Project struct {
	Type                           string              `json:"type"                       mapstructure:"type"`
	ID                             common.UUID         `json:"uuid"                       mapstructure:"uuid"`
	Name                           string              `json:"name"                       mapstructure:"name"`
	Description                    string              `json:"description,omitempty"      mapstructure:"description"`
	Key                            string              `json:"key"                        mapstructure:"key"`
	Owner                          user.User           `json:"owner"                      mapstructure:"owner"`
	Workspace                      workspace.Workspace `json:"workspace"                  mapstructure:"workspace"`
	Links                          common.Links        `json:"links"                      mapstructure:"links"`
	IsPrivate                      bool                `json:"is_private"                 mapstructure:"is_private"`
	HasPubliclyVisibleRepositories bool                `json:"has_publicly_visible_repos" mapstructure:"has_publicly_visible_repos"`
	CreatedOn                      time.Time           `json:"created_on"                 mapstructure:"created_on"`
	UpdatedOn                      time.Time           `json:"updated_on"                 mapstructure:"updated_on"`
}

type ProjectReference struct {
	Key string `json:"key" mapstructure:"key"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = []string{
	"key",
	"name",
	"description",
	"owner",
	"workspace",
	"created_on",
	"updated_on",
	"private",
}

func init() {
	Command.AddCommand(reviewer.Command)
}

// NewReference creates a new ProjectReference
func NewReference(key string) *ProjectReference {
	return &ProjectReference{
		Key: key,
	}
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (project Project) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Key", "Name", "Description"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (project Project) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "key":
			row = append(row, project.Key)
		case "name":
			row = append(row, project.Name)
		case "description":
			row = append(row, project.Description)
		case "owner":
			if project.Owner.Name == "" {
				row = append(row, " ")
			} else {
				row = append(row, project.Owner.Name)
			}
		case "workspace":
			if project.Workspace.Name == "" {
				row = append(row, " ")
			} else {
				row = append(row, project.Workspace.Name)
			}
		case "created on", "created-on", "created_on", "created":
			row = append(row, project.CreatedOn.Format("2006-01-02 15:04:05"))
		case "updated on", "updated-on", "updated_on", "updated":
			if !project.UpdatedOn.IsZero() {
				row = append(row, project.UpdatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		case "private":
			row = append(row, fmt.Sprintf("%t", project.IsPrivate))
		}
	}
	return row
}

// Validate validates a Project
func (project *Project) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pullrequest
//
// implements fmt.Stringer
func (project Project) String() string {
	return project.Name
}

// MarshalJSON implements the json.Marshaler interface.
func (project Project) MarshalJSON() (data []byte, err error) {
	type surrogate Project
	var owner *user.User
	var wspace *workspace.Workspace
	var createdOn string
	var updatedOn string

	if !project.Owner.ID.IsNil() {
		owner = &project.Owner
	}
	if !project.Workspace.ID.IsNil() {
		wspace = &project.Workspace
	}
	if !project.CreatedOn.IsZero() {
		createdOn = project.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}
	if !project.UpdatedOn.IsZero() {
		updatedOn = project.UpdatedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}

	data, err = json.Marshal(struct {
		surrogate
		Owner     *user.User           `json:"owner,omitempty"`
		Workspace *workspace.Workspace `json:"workspace,omitempty"`
		CreatedOn string               `json:"created_on,omitempty"`
		UpdatedOn string               `json:"updated_on,omitempty"`
	}{
		surrogate: surrogate(project),
		Owner:     owner,
		Workspace: wspace,
		CreatedOn: createdOn,
		UpdatedOn: updatedOn,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// GetWorkspace gets the workspace from the command
func GetWorkspace(cmd *cobra.Command, profile *profile.Profile) (workspace string, err error) {
	workspace = cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = profile.DefaultWorkspace
		if len(workspace) == 0 {
			return "", errors.ArgumentMissing.With("workspace")
		}
	}
	return
}

// GetProjectKeys gets the keys of the projects in the workspace given in the command
func GetProjectKeys(context context.Context, cmd *cobra.Command, args []string, toComplete string) (keys []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = profile.Current.DefaultWorkspace
		if len(workspace) == 0 {
			log.Warnf("No workspace given")
			return
		}
	}

	projects, err := profile.GetAll[Project](context, cmd, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	keys = core.Map(projects, func(project Project) string { return project.Key })
	core.Sort(keys, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return keys, nil
}

// GetProjectNames gets the names of the projects in the workspace given in the command
func GetProjectNames(context context.Context, cmd *cobra.Command, args []string, toComplete string) (names []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("project", "names")

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = profile.Current.DefaultWorkspace
		if len(workspace) == 0 {
			log.Warnf("No workspace given")
			return
		}
	}

	log.Infof("Getting all projects from workspace %s", workspace)
	projects, err := profile.GetAll[Project](context, cmd, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	names = core.Map(projects, func(project Project) string { return project.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return names, nil
}
