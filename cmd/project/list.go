package project

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all projects",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Workspace common.RemoteValueFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	listCmd.Flags().Var(&listOptions.Workspace, "workspace", "Workspace to list projects from")
	_ = listCmd.MarkFlagRequired("workspace")
	_ = listCmd.RegisterFlagCompletionFunc("workspace", listOptions.Workspace.CompletionFunc())
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all projects from workspace %s with profile %s", listOptions.Workspace, profile.Current)
	projects, err := profile.GetAll[Project](
		cmd.Context(),
		cmd,
		profile.Current,
		fmt.Sprintf("/workspaces/%s/projects", listOptions.Workspace),
	)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		log.Infof("No project found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), Projects(projects))
}
