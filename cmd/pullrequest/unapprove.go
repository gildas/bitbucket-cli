package pullrequest

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unapproveCmd = &cobra.Command{
	Use:               "unapprove [flags] <pullrequest-id>",
	Short:             "unapprove a pullrequest by its <pullrequest-id>. If not provided, it will try to unapprove the only open pullrequest.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: unapproveValidArgs,
	RunE:              unapproveProcess,
}

func init() {
	Command.AddCommand(unapproveCmd)
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

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot unapprove Pull Request"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot unapprove Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot unapprove Pull Request"), err)
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unapproving pullrequest %s", pullRequestID) {
		return nil
	}
	err = profile.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", pullRequestID, "approve"),
		nil,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to unapprove Pull Request %s", pullRequestID), err)
	}
	return
}
