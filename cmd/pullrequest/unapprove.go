package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unapproveCmd = &cobra.Command{
	Use:               "unapprove [flags] <pullrequest-id>",
	Short:             "unapprove a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: unapproveValidArgs,
	RunE:              unapproveProcess,
}

var unapproveOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(unapproveCmd)

	unapproveCmd.Flags().StringVar(&unapproveOptions.Repository, "repository", "", "Repository to unapprove pullrequest from. Defaults to the current repository")
}

func unapproveValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	ids, err := prcommon.GetPullRequestIDsWithState(cmd.Context(), cmd, "OPEN")
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func unapproveProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unapprove")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unapproving pullrequest %s", args[0]) {
		return nil
	}
	err = profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/approve", args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unapprove pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return
}
