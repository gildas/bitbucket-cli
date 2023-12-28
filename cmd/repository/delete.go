package repository

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <slug_or_uuid>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete a repository by its <slug> or <uuid>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Workspace common.RemoteValueFlag
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	deleteCmd.Flags().Var(&deleteOptions.Workspace, "workspace", "Workspace to delete repositories from")
	_ = deleteCmd.MarkFlagRequired("workspace")
	_ = deleteCmd.RegisterFlagCompletionFunc("workspace", deleteOptions.Workspace.CompletionFunc())
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetRepositorySlugs(cmd.Context(), cmd, profile.Current, deleteOptions.Workspace.String()), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Deleting repository %s", args[0])

	err := profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/repositories/%s/%s", deleteOptions.Workspace, args[0]),
		nil,
	)
	if err != nil {
		return err
	}

	log.Infof("Repository %s deleted", args[0])
	return nil
}
