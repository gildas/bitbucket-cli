package workspace

import (
	"context"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Workspace struct {
	Type  string     `json:"type"  mapstructure:"type"`
	ID    string     `json:"uuid"  mapstructure:"uuid"`
	Name  string     `json:"name"  mapstructure:"name"`
	Slug  string     `json:"slug"  mapstructure:"slug"`
	Links link.Links `json:"links" mapstructure:"links"`
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
		workspace.ID,
		workspace.Name,
		workspace.Slug,
	}
}

// GetWorkspaceSlugs gets the slugs of all workspaces
func GetWorkspaceSlugs(context context.Context) (slugs []string) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "slugs")

	log.Debugf("Getting all workspaces")
	workspaces, err := profile.GetAll[Workspace](
		context,
		profile.Current,
		"",
		"/workspaces",
	)
	if err != nil {
		log.Errorf("Failed to get workspaces", err)
		return
	}
	return core.Map(workspaces, func(workspace Workspace) string {
		return workspace.Slug
	})
}
