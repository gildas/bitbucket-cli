package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
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
	_ = getCmd.RegisterFlagCompletionFunc("project", deleteOptions.Project.CompletionFunc("project"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	workspace := deleteOptions.Workspace.Value
	if len(workspace) == 0 {
		workspace = profile.Current.DefaultWorkspace
		if len(workspace) == 0 {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return GetReviewerUserIDs(cmd.Context(), cmd, profile.Current, workspace, deleteOptions.Project.Value), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(deleteOptions.Workspace.Value) == 0 {
		deleteOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(deleteOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}
	if len(deleteOptions.Project.Value) == 0 {
		deleteOptions.Project.Value = profile.Current.DefaultProject
		if len(deleteOptions.Project.Value) == 0 {
			return errors.ArgumentMissing.With("project")
		}
	}

	var merr errors.MultiError
	for _, userID := range args {
		if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting default reviewer %s from project %s", userID, deleteOptions.Project) {
			err := profile.Current.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", deleteOptions.Workspace, deleteOptions.Project, userID),
				nil,
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete default reviewer %s: %s\n", userID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Default reviewer %s deleted", userID)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these reviewers: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these reviewers, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
