package comment

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	prcommon "bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pullrequest comments",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository    string
	PullRequestID *flags.EnumFlag
	Columns       *flags.EnumSliceFlag
	SortBy        *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequest comments from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.PullRequestID, "pullrequest", "pullrequest to list comments from")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.MarkFlagRequired("pullrequest")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all comments from repository %s with profile %s", listOptions.Repository, profile.Current)
	comments, err := profile.GetAll[Comment](
		cmd.Context(),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments", listOptions.PullRequestID.Value),
	)
	if err != nil {
		return err
	}
	if len(comments) == 0 {
		log.Infof("No comment found")
		return nil
	}
	core.Sort(comments, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(
		cmd.Context(),
		cmd,
		Comments(core.Filter(comments, func(comment Comment) bool {
			return len(comment.Content.Raw) > 0
		})),
	)
}
