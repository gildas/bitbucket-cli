package project

import (
	"context"
	"encoding/json"
	"fmt"
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

func init() {
	Command.AddCommand(reviewer.Command)
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (project Project) GetHeader(short bool) []string {
	return []string{"Key", "Name", "Description"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (project Project) GetRow(headers []string) []string {
	return []string{project.Key, project.Name, project.Description}
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

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
	}{
		surrogate: surrogate(project),
		CreatedOn: project.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
		UpdatedOn: project.UpdatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// GetProjectKeys gets the keys of the projects in the workspace given in the command
func GetProjectKeys(context context.Context, cmd *cobra.Command) (keys []string) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		log.Warnf("No workspace given")
		return
	}

	projects, err := profile.GetAll[Project](context, cmd, profile.Current, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	return core.Map(projects, func(project Project) string {
		return project.Key
	})
}

// GetProjectNames gets the names of the projects in the workspace given in the command
func GetProjectNames(context context.Context, cmd *cobra.Command) (keys []string) {
	log := logger.Must(logger.FromContext(context)).Child("project", "names")

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		log.Warnf("No workspace given")
		return
	}

	log.Infof("Getting all projects from workspace %s", workspace)
	projects, err := profile.GetAll[Project](context, cmd, profile.Current, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	return core.Map(projects, func(project Project) string {
		return project.Name
	})
}
