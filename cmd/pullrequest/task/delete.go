package task

import (
	"fmt"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <task-id...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete pullrequest tasks by their <task-id>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	PullRequestID *flags.EnumFlag
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.PullRequestID = flags.NewEnumFlagWithFunc(deleteCmd, "", prcommon.GetPullRequestIDs)
	deleteCmd.Flags().Var(deleteOptions.PullRequestID, "pullrequest", "Pullrequest to delete comments from")
	_ = deleteCmd.MarkFlagRequired("pullrequest")
	_ = deleteCmd.RegisterFlagCompletionFunc(deleteOptions.PullRequestID.CompletionFunc("pullrequest"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	taskIDs, err := GetPullRequestTaskIDs(cmd.Context(), cmd, deleteOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return taskIDs, cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, taskID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting task %s from pullrequest %s", taskID, deleteOptions.PullRequestID) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				repository.GetPath("pullrequests", deleteOptions.PullRequestID.Value, "tasks", taskID),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete pullrequest task %s: %s\n", taskID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Pullrequest task %s deleted", taskID)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these tasks: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these tasks, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
