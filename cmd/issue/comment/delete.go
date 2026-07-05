package comment

import (
	"fmt"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <comment-id...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete issue comments by their <comment-id>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	IssueID *flags.EnumFlag
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.IssueID = flags.NewEnumFlagWithFunc(deleteCmd, "", GetIssueIDs)
	deleteCmd.Flags().Var(deleteOptions.IssueID, "issue", "Issue to delete comments from")
	_ = deleteCmd.MarkFlagRequired("issue")
	_ = deleteCmd.RegisterFlagCompletionFunc(deleteOptions.IssueID.CompletionFunc("issue"))
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	commentIDs, err := GetIssueCommentIDs(cmd.Context(), cmd, profile.Current, deleteOptions.IssueID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(commentIDs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
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
	for _, commentID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting comment %s from issue %s", commentID, deleteOptions.IssueID) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				repository.GetPath("issues", deleteOptions.IssueID.Value, "comments", commentID),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					return errors.Join(errors.Errorf("Failed to delete issue comment %s", commentID), err)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Issue comment %s deleted", commentID)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these comments: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these comments, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
