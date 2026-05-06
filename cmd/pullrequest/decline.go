package pullrequest

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/bitbucket-cli/cmd/user"
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
		return errors.Join(errors.Errorf("Cannot decline Pull Request"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot decline Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot decline Pull Request"), err)
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
		return errors.Join(errors.Errorf("Failed to decline Pull Request %s", pullRequestID), err)
	}
	return profile.Print(cmd.Context(), cmd, participant)
}
