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
	Message           string
	MergeStrategy     *flags.EnumFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(mergeCmd)

	mergeOptions.MergeStrategy = flags.NewEnumFlag("+merge_commit", "squash", "fast_forward")
	mergeCmd.Flags().StringVar(&mergeOptions.Message, "message", "", "Message of the merge")
	mergeCmd.Flags().BoolVar(&mergeOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch of the pullrequest")
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

	var pullrequest PullRequest

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
	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", pullRequestID, "merge"),
		payload,
		&pullrequest,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to merge Pull Request %s", pullRequestID), err)
	}
	return profile.Print(cmd.Context(), cmd, pullrequest)
}
