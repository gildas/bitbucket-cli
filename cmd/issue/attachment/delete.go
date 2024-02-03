package attachment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <path...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete issue attachments by their <path>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	IssueID      *flags.EnumFlag
	Repository   string
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete an issue attachment from. Defaults to the current repository")
	deleteCmd.Flags().Var(deleteOptions.IssueID, "issue", "Issue to delete attachments from")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
	_ = deleteCmd.MarkFlagRequired("issue")
	_ = deleteCmd.RegisterFlagCompletionFunc("issue", deleteOptions.IssueID.CompletionFunc("issue"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetAttachmentNames(cmd.Context(), cmd, profile.Current, downloadOptions.IssueID.Value), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var merr errors.MultiError
	for _, attachmentID := range args {
		if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting attachment %s from issue %s", attachmentID, deleteOptions.IssueID) {
			err := profile.Current.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("issues/%s/attachments/%s", deleteOptions.IssueID.Value, attachmentID),
				nil,
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete issue attachment %s: %s\n", attachmentID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Issue attachment %s deleted", attachmentID)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these attachments: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these attachments, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
