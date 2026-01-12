package repository

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <slug_or_uuid...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete repositories by their <slug> or <uuid>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Workspace *flags.EnumFlag
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	deleteCmd.Flags().Var(deleteOptions.Workspace, "workspace", "Workspace to delete repositories from")
	_ = deleteCmd.RegisterFlagCompletionFunc(deleteOptions.Workspace.CompletionFunc("workspace"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	slugs, err := GetRepositorySlugs(cmd.Context(), cmd, deleteOptions.Workspace.String())
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(slugs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if len(deleteOptions.Workspace.Value) == 0 {
		deleteOptions.Workspace.Value = profile.DefaultWorkspace
		if len(deleteOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	var merr errors.MultiError
	for _, repositorySlug := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting repository %s", repositorySlug) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("/repositories/%s/%s", deleteOptions.Workspace, repositorySlug),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete repository %s: %s\n", repositorySlug, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
		}
		log.Infof("Repository %s deleted", repositorySlug)
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these repositories: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these repositories, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
