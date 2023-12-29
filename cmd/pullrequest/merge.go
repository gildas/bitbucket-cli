package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:               "merge [flags] <pullrequest-id>",
	Short:             "merge a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: mergeValidArgs,
	RunE:              mergeProcess,
}

var mergeOptions struct {
	Repository        string
	Message           string
	MergeStrategy     common.EnumFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(mergeCmd)

	mergeOptions.MergeStrategy = common.EnumFlag{Allowed: []string{"merge_commit", "squash", "fast_forward"}, Value: "merge_commit"}
	mergeCmd.Flags().StringVar(&mergeOptions.Repository, "repository", "", "Repository to merge pullrequest from. Defaults to the current repository")
	mergeCmd.Flags().StringVar(&mergeOptions.Message, "message", "", "Message of the merge")
	mergeCmd.Flags().BoolVar(&mergeOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch of the pullrequest")
	mergeCmd.Flags().Var(&mergeOptions.MergeStrategy, "merge-strategy", "Merge strategy to use. Possible values are \"merge_commit\", \"squash\" or \"fast_forward\"")
	_ = mergeCmd.RegisterFlagCompletionFunc("merge-strategy", mergeOptions.MergeStrategy.CompletionFunc())
}

func mergeValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return GetPullRequestIDs(cmd.Context(), cmd, mergeOptions.Repository, "OPEN"), cobra.ShellCompDirectiveNoFileComp
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

	log.Record("payload", payload).Infof("Merging pullrequest %s", args[0])
	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/merge", args[0]),
		payload,
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to merge pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), pullrequest)
}
