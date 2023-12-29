package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a reviewer",
	ValidArgsFunction: getValidArgs,
	Args:              cobra.ExactArgs(1),
	RunE:              getProcess,
}

var getOptions struct {
	Workspace common.RemoteValueFlag
	Project   common.RemoteValueFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	getOptions.Project = common.RemoteValueFlag{AllowedFunc: GetProjectKeys}
	getCmd.Flags().Var(&getOptions.Workspace, "workspace", "Workspace to get reviewers from")
	getCmd.Flags().Var(&getOptions.Project, "project", "Project Key to get reviewers from")
	_ = getCmd.RegisterFlagCompletionFunc("workspace", getOptions.Workspace.CompletionFunc())
	_ = getCmd.RegisterFlagCompletionFunc("project", getOptions.Project.CompletionFunc())
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	workspace := getOptions.Workspace.Value
	if len(workspace) == 0 {
		workspace = profile.Current.DefaultWorkspace
		if len(workspace) == 0 {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return GetReviewerUserIDs(cmd.Context(), cmd, profile.Current, workspace, getOptions.Project.Value), cobra.ShellCompDirectiveNoFileComp
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
	if len(getOptions.Project.Value) == 0 {
		getOptions.Project.Value = profile.Current.DefaultProject
		if len(getOptions.Project.Value) == 0 {
			return errors.ArgumentMissing.With("project")
		}
	}

	log.Infof("Displaying reviewer %s", args[0])
	var user user.User

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", getOptions.Workspace, getOptions.Project, args[0]),
		&user,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get reviewer: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), user)
}
