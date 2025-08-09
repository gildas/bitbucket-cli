package project

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
	Use:               "delete [flags] <project-key...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete projects by their <project-key>.",
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
	deleteCmd.Flags().Var(deleteOptions.Workspace, "workspace", "Workspace to delete projects from")
	_ = deleteCmd.RegisterFlagCompletionFunc(createOptions.Workspace.CompletionFunc("workspace"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	keys, err := GetProjectKeys(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return common.FilterValidArgs(keys, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, err := GetWorkspace(cmd, profile)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, projectKey := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting project %s", projectKey) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("/workspaces/%s/projects/%s", workspace, projectKey),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete project %s: %s\n", projectKey, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Project %s deleted", projectKey)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these projects: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these projects, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
