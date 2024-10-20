package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	errors "github.com/gildas/go-errors"
	flags "github.com/gildas/go-flags"
	logger "github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <user-id...>",
	Aliases:           []string{"remove"},
	Short:             "delete  reviewers by their <user-id>.",
	ValidArgsFunction: deleteValidArgs,
	Args:              cobra.MinimumNArgs(1),
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Workspace    *flags.EnumFlag
	Project      *flags.EnumFlag
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	deleteOptions.Project = flags.NewEnumFlagWithFunc("", GetProjectKeys)
	deleteCmd.Flags().Var(deleteOptions.Workspace, "workspace", "Workspace to delete reviewers from")
	deleteCmd.Flags().Var(deleteOptions.Project, "project", "Project Key to delete reviewers from")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
	_ = deleteCmd.RegisterFlagCompletionFunc("workspace", deleteOptions.Workspace.CompletionFunc("workspace"))
	_ = deleteCmd.RegisterFlagCompletionFunc("project", deleteOptions.Project.CompletionFunc("project"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	userIDs, err := GetReviewerUserIDs(cmd.Context(), cmd, deleteOptions.Project.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return userIDs, cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, project, err := GetWorkspaceAndProject(cmd, profile)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, userID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting default reviewer %s from project %s", userID, project) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", workspace, project, userID),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete default reviewer %s: %s\n", userID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Default reviewer %s deleted", userID)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these reviewers: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these reviewers, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
