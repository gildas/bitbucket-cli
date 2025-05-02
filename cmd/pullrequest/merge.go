package pullrequest

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:               "merge [flags] <pullrequest-id>",
	Short:             "merge a pullrequest by its <pullrequest-id>.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: mergeValidArgs,
	RunE:              mergeProcess,
}

var mergeOptions struct {
	Repository        string
	Message           string
	MergeStrategy     *flags.EnumFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(mergeCmd)

	mergeOptions.MergeStrategy = flags.NewEnumFlag("+merge_commit", "squash", "fast_forward")
	mergeCmd.Flags().StringVar(&mergeOptions.Repository, "repository", "", "Repository to merge pullrequest from. Defaults to the current repository")
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

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
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

	var pullRequestID string

	if len(args) == 0 {
		pullRequestIDs, err := prcommon.GetPullRequestIDsWithState(cmd.Context(), cmd, "OPEN")
		if err != nil {
			return err
		}
		if len(pullRequestIDs) == 0 {
			return errors.Errorf("No pullrequest to merge")
		}
		if len(pullRequestIDs) > 1 {
			return errors.Errorf("Too many pullrequests to merge: %s", strings.Join(pullRequestIDs, ", "))
		}
		pullRequestID = pullRequestIDs[0]
	} else {
		pullRequestID = args[0]
	}

	if _, err := strconv.Atoi(pullRequestID); err != nil {
		return errors.ArgumentInvalid.With("pullrequest-id", pullRequestID)
	}

	log.Record("payload", payload).Infof("Merging pullrequest %s", pullRequestID)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Merging pullrequest %s", pullRequestID) {
		return nil
	}
	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/merge", pullRequestID),
		payload,
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to merge pullrequest %s: %s\n", pullRequestID, err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, pullrequest)
}
