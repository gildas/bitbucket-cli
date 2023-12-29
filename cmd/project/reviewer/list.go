package reviewer

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
	Short: "list all reviewers",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Workspace common.RemoteValueFlag
	Project   common.RemoteValueFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	listOptions.Project = common.RemoteValueFlag{AllowedFunc: GetProjectKeys}
	listCmd.Flags().Var(&listOptions.Workspace, "workspace", "Workspace to list reviewers from")
	listCmd.Flags().Var(&listOptions.Project, "project", "Project Key to list reviewers from")
	_ = listCmd.RegisterFlagCompletionFunc("workspace", listOptions.Workspace.CompletionFunc())
	_ = getCmd.RegisterFlagCompletionFunc("project", listOptions.Project.CompletionFunc())
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(listOptions.Workspace.Value) == 0 {
		listOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(listOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}
	if len(listOptions.Project.Value) == 0 {
		listOptions.Project.Value = profile.Current.DefaultProject
		if len(listOptions.Project.Value) == 0 {
			return errors.ArgumentMissing.With("project")
		}
	}

	log.Infof("Listing all reviewers")
	reviewers, err := profile.GetAll[Reviewer](
		cmd.Context(),
		cmd,
		profile.Current,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers", listOptions.Workspace, listOptions.Project),
	)
	if err != nil {
		return err
	}
	if len(reviewers) == 0 {
		log.Infof("No reviewer found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), Reviewers(reviewers))
}
