package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var removeRequestChangesCmd = &cobra.Command{
	Use:               "remove-request-changes [flags] <pullrequest-id>",
	Aliases:           []string{"removeRequestChanges", "remove-requestChanges", "removerequestchanges", "cancel-request-changes"},
	Short:             "Remove request changes on a pullrequest by its <pullrequest-id>. If not provided, it will try to remove request changes on the only open pullrequest.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: removeRequestChangesValidArgs,
	RunE:              removeRequestChangesProcess,
}

func init() {
	Command.AddCommand(removeRequestChangesCmd)
}

func removeRequestChangesValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")
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
	log.Debugf("Fetched %d pullrequest ids", len(ids))
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func removeRequestChangesProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "removeRequestChanges")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot remove request changes on Pull Request"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot remove request changes on Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot remove request changes on Pull Request"), err)
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Removing request changes on pullrequest %s", pullRequestID) {
		return nil
	}

	err = profile.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", pullRequestID, "request-changes"),
		nil,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to remove request changes on Pull Request %s", pullRequestID), err)
	}
	return
}
