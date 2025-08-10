package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <comment-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a pullrequest comment by its <comment-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	PullRequestID *flags.EnumFlag
	Repository    string
	Columns       *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
	getOptions.Columns = flags.NewEnumSliceFlag(columns...)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get a pullrequest comment from. Defaults to the current repository")
	getCmd.Flags().Var(getOptions.PullRequestID, "pullrequest", "Pullrequest to get comments from")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.MarkFlagRequired("pullrequest")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	commentIDs, err := GetPullRequestCommentIDs(cmd.Context(), cmd, deleteOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return commentIDs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying pullrequest comment %s", args[0])
	var comment Comment

	err = profile.Current.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments/%s", getOptions.PullRequestID.Value, args[0]),
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pullrequest comment %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, comment)
}
