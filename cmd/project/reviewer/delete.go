package reviewer

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

var deleteCmd = &cobra.Command{
	Use:               "delete",
	Aliases:           []string{"remove"},
	Short:             "delete a reviewer",
	ValidArgsFunction: deleteValidArgs,
	Args:              cobra.ExactArgs(1),
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Workspace common.RemoteValueFlag
	Project   common.RemoteValueFlag
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	deleteOptions.Project = common.RemoteValueFlag{AllowedFunc: GetProjectKeys}
	deleteCmd.Flags().Var(&deleteOptions.Workspace, "workspace", "Workspace to delete reviewers from")
	deleteCmd.Flags().Var(&deleteOptions.Project, "project", "Project Key to delete reviewers from")
	_ = deleteCmd.MarkFlagRequired("workspace")
	_ = deleteCmd.MarkFlagRequired("project")
	_ = deleteCmd.RegisterFlagCompletionFunc("workspace", deleteOptions.Workspace.CompletionFunc())
	_ = getCmd.RegisterFlagCompletionFunc("project", deleteOptions.Project.CompletionFunc())
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetReviewerUserIDs(cmd.Context(), cmd, profile.Current, deleteOptions.Workspace.Value, deleteOptions.Project.Value), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("deleteing reviewer %s", args[0])

	err := profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", deleteOptions.Workspace, deleteOptions.Project, args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete reviewer: %s\n", err)
		os.Exit(1)
	}
	return nil
}
