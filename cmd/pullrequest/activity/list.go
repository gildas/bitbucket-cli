package activity

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pullrequest Activities",
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
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns...)
	listOptions.SortBy = flags.NewEnumFlag(sortBy...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequest activities from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.PullRequestID, "pullrequest", "pullrequest to list activities from")
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

	log.Infof("Listing all activities from repository %s with profile %s", listOptions.Repository, profile.Current)
	activities, err := profile.GetAll[Activity](
		cmd.Context(),
		cmd,
		fmt.Sprintf("pullrequests/%s/activity", listOptions.PullRequestID.Value),
	)
	if err != nil {
		return err
	}
	if len(activities) == 0 {
		log.Infof("No activities found")
		return nil
	}
	core.Sort(activities, func(a, b Activity) bool {
		if a.PullRequest.ID < b.PullRequest.ID {
			return true
		}
		if a.Update != nil && b.Update != nil {
			return a.Update.Date.Before(b.Update.Date)
		}
		if a.Approval != nil && b.Approval != nil {
			return a.Approval.Date.Before(b.Approval.Date)
		}
		return false
	})
	return profile.Current.Print(
		cmd.Context(),
		cmd,
		Activities(core.Filter(activities, func(activity Activity) bool {
			return true
		})),
	)
}
