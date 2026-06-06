package pullrequest

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var mergeStatusCmd = &cobra.Command{
	Use:               "merge-status <pull-request-id>",
	Short:             "Get the status of a pull request merge task",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: mergeStatusValidArgs,
	RunE:              mergeStatusProcess,
}

var mergeStatusOptions struct {
	TaskID string
}

func init() {
	Command.AddCommand(mergeStatusCmd)

	mergeStatusCmd.Flags().StringVar(&mergeStatusOptions.TaskID, "task-id", "", "ID of the merge task to check the status of")
	mergeStatusCmd.MarkFlagRequired("task-id")
}

func mergeStatusValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ids, err := prcommon.GetPullRequestIDsWithState(cmd.Context(), cmd, "OPEN")
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func mergeStatusProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "merge-status")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to get the profile"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	log.Infof("Getting the Pull Request merge status for %s", pullRequestID)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Getting the merge status for pull request %s", pullRequestID) {
		return nil
	}

	var status PullRequestMergeStatus

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", pullRequestID, "merge", "task-status", mergeStatusOptions.TaskID),
		&status,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to get the merge status for Pull Request %s", pullRequestID), err)
	}
	status.ID = mergeStatusOptions.TaskID

	return profile.Print(cmd.Context(), cmd, status)
}
