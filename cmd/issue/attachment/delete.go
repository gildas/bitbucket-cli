package attachment

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
	Use:               "delete [flags] <path...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete issue attachments by their <path>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	IssueID    *flags.EnumFlag
	Repository string
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete an issue attachment from. Defaults to the current repository")
	deleteCmd.Flags().Var(deleteOptions.IssueID, "issue", "Issue to delete attachments from")
	_ = deleteCmd.MarkFlagRequired("issue")
	_ = deleteCmd.RegisterFlagCompletionFunc(deleteOptions.IssueID.CompletionFunc("issue"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	attachmentNames, err := GetAttachmentNames(cmd.Context(), cmd, deleteOptions.IssueID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(attachmentNames, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, attachmentID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting attachment %s from issue %s", attachmentID, deleteOptions.IssueID) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("issues/%s/attachments/%s", deleteOptions.IssueID.Value, attachmentID),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete issue attachment %s: %s\n", attachmentID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Issue attachment %s deleted", attachmentID)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these attachments: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these attachments, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
