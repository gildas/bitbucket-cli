package pullrequest

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var declineCmd = &cobra.Command{
	Use:               "decline [flags] <pullrequest-id>",
	Short:             "decline a pullrequest by its <pullrequest-id>. If not provided, it will try to decline the only open pullrequest.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: declineValidArgs,
	RunE:              declineProcess,
}

func init() {
	Command.AddCommand(declineCmd)
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

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var pullRequestID string

	if len(args) == 0 {
		pullRequestIDs, err := prcommon.GetPullRequestIDsFromRepositoryWithState(cmd.Context(), cmd, repository, "OPEN")
		if err != nil {
			return err
		}
		if len(pullRequestIDs) == 0 {
			return errors.Errorf("No pullrequest to decline")
		}
		if len(pullRequestIDs) > 1 {
			return errors.Errorf("Too many pullrequests to decline: %s", strings.Join(pullRequestIDs, ", "))
		}
		pullRequestID = pullRequestIDs[0]
	} else {
		pullRequestID = args[0]
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Declining pullrequest %s", pullRequestID) {
		return nil
	}
	var participant user.Participant

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", pullRequestID, "decline"),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decline pullrequest %s: %s\n", pullRequestID, err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, participant)
}
