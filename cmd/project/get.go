package project

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <project-key>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a project by its <project-key>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Workspace common.RemoteValueFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	getCmd.Flags().Var(&getOptions.Workspace, "workspace", "Workspace to get projects from")
	_ = getCmd.RegisterFlagCompletionFunc("workspace", getOptions.Workspace.CompletionFunc())
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetProjectKeys(cmd.Context(), cmd, args), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(getOptions.Workspace.Value) == 0 {
		getOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(getOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	log.Infof("Displaying project %s", args[0])
	var project Project

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s", getOptions.Workspace, args[0]),
		&project,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, project)
}
