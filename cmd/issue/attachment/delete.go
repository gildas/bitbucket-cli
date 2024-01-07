package attachment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <path>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete an issue attachment by its <path>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	IssueID    common.RemoteValueFlag
	Repository string
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteOptions.IssueID = common.RemoteValueFlag{AllowedFunc: GetIssueIDs}
	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete an issue attachment from. Defaults to the current repository")
	deleteCmd.Flags().Var(&deleteOptions.IssueID, "issue", "Issue to delete attachments from")
	_ = deleteCmd.MarkFlagRequired("issue")
	_ = deleteCmd.RegisterFlagCompletionFunc("issue", deleteOptions.IssueID.CompletionFunc())
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

	if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting attachment %s from issue %s", args[0], deleteOptions.IssueID) {
		err := profile.Current.Delete(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/attachments/%s", deleteOptions.IssueID.Value, args[0]),
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete issue artifact %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return nil
}
