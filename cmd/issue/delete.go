package issue

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
	Use:               "delete [flags] <issue-id...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete issues by their <issue-id>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Repository   string
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete an issue from. Defaults to the current repository")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ids, err := GetIssueIDs(cmd.Context(), cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, issueID := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting issue %s", issueID) {
			err := profile.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("issues/%s", issueID),
				nil,
			)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete issue %s: %s\n", issueID, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Issue %s deleted", issueID)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these issues: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these issues, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
