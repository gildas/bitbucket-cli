package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var declineCmd = &cobra.Command{
	Use:               "decline [flags] <pullrequest-id>",
	Short:             "decline a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: declineValidArgs,
	RunE:              declineProcess,
}

var declineOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(declineCmd)

	declineCmd.Flags().StringVar(&declineOptions.Repository, "repository", "", "Repository to decline pullrequest from. Defaults to the current repository")
}

func declineValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func declineProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "decline")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Declining pullrequest %s", args[0]) {
		return nil
	}
	var participant user.Participant

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/decline", args[0]),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decline pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, participant)
}
