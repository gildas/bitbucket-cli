package common

import (
	"context"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
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
