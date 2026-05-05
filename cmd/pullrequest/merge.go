package pullrequest

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:               "merge [flags] <pullrequest-id>",
	Short:             "merge a pullrequest by its <pullrequest-id>. If not provided, it will try to merge the only open pullrequest.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: mergeValidArgs,
	RunE:              mergeProcess,
}

var mergeOptions struct {
	Async             bool
	Message           string
	MergeStrategy     *flags.EnumFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(mergeCmd)

	mergeOptions.MergeStrategy = flags.NewEnumFlag("+merge_commit", "squash", "fast_forward")
	mergeCmd.Flags().StringVar(&mergeOptions.Message, "message", "", "Message of the merge")
	mergeCmd.Flags().BoolVar(&mergeOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch of the pullrequest")
	mergeCmd.Flags().BoolVar(&mergeOptions.Async, "async", false, "Perform the merge asynchronously")
	mergeCmd.Flags().Var(mergeOptions.MergeStrategy, "merge-strategy", "Merge strategy to use. Possible values are \"merge_commit\", \"squash\" or \"fast_forward\"")
	_ = mergeCmd.RegisterFlagCompletionFunc(mergeOptions.MergeStrategy.CompletionFunc("merge-strategy"))
}

func mergeValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func mergeProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "merge")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	uripath := repository.GetPath("pullrequests", pullRequestID, "merge")

	if mergeOptions.Async {
		uripath += "?async=true"
	}

	payload := struct {
		Message           string `json:"message,omitempty"`
		CloseSourceBranch bool   `json:"close_source_branch"`
		MergeStrategy     string `json:"merge_strategy"`
	}{
		Message:           mergeOptions.Message,
		CloseSourceBranch: mergeOptions.CloseSourceBranch,
		MergeStrategy:     mergeOptions.MergeStrategy.String(),
	}

	log.Record("payload", payload).Infof("Merging pullrequest %s", pullRequestID)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Merging pullrequest %s", pullRequestID) {
		return nil
	}

	if mergeOptions.Async {
		result, err := profile.PostWithResult(log.ToContext(cmd.Context()), cmd, uripath, payload)
		if err != nil {
			return errors.Join(errors.Errorf("Failed to merge Pull Request %s", pullRequestID), err)
		}
		status, err := NewPullRequestMergeStatusFromLocation(result.Headers.Get("Location"))
		if err != nil {
			return errors.Join(errors.Errorf("Failed to get merge status for Pull Request %s", pullRequestID), err)
		}
		log.Infof("Merge request accepted, task ID: %s", status.ID)
		return profile.Print(cmd.Context(), cmd, status)
	} else {
		var pullrequest PullRequest

		err = profile.Post(log.ToContext(cmd.Context()), cmd, uripath, payload, &pullrequest)
		if err != nil {
			return errors.Join(errors.Errorf("Failed to merge Pull Request %s", pullRequestID), err)
		}
		return profile.Print(cmd.Context(), cmd, pullrequest)
	}
}
