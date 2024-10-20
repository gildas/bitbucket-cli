package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <comment-id...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete pullrequest comments by their <comment-id>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	PullRequestID *flags.EnumFlag
	Repository    string
	StopOnError   bool
	WarnOnError   bool
	IgnoreErrors  bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.PullRequestID = flags.NewEnumFlagWithFunc("", GetPullRequestIDs)
	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete a pullrequest comment from. Defaults to the current repository")
	deleteCmd.Flags().Var(deleteOptions.PullRequestID, "pullrequest", "Pullrequest to delete comments from")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	_ = deleteCmd.MarkFlagRequired("pullrequest")
	_ = deleteCmd.RegisterFlagCompletionFunc("pullrequest", deleteOptions.PullRequestID.CompletionFunc("pullrequest"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	commentIDs, err := GetPullRequestCommentIDs(cmd.Context(), cmd, deleteOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return commentIDs, cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var merr errors.MultiError
	for _, commentID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting comment %s from pullrequest %s", commentID, deleteOptions.PullRequestID) {
			err := profile.Current.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("pullrequests/%s/comments/%s", deleteOptions.PullRequestID.Value, commentID),
				nil,
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete pullrequest comment %s: %s\n", commentID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Pullrequest comment %s deleted", commentID)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these comments: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these comments, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
