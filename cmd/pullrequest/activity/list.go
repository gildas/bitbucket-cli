package activity

import (
	"fmt"
	"net/url"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:        "list",
	Short:      "list all pullrequest Activities",
	Deprecated: "Please use 'pullrequest activities' instead",
	Args:       cobra.MaximumNArgs(1),
	RunE:       listProcess,
}

var listOptions struct {
	PullRequestID *flags.EnumFlag
	Query         string
	Columns       *flags.EnumSliceFlag
	SortBy        *flags.EnumFlag
	PageLength    int
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.PullRequestID = flags.NewEnumFlagWithFunc(listCmd, "", prcommon.GetPullRequestIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().Var(listOptions.PullRequestID, "pullrequest", "pullrequest to list activities from")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter activities")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	listCmd.Flags().IntVar(&listOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
	_ = listCmd.MarkFlagRequired("pullrequest")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	uripath := repository.GetPath(fmt.Sprintf("pullrequests/%s/activity", listOptions.PullRequestID.Value))

	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("%s?q=%s", uripath, url.QueryEscape(listOptions.Query))
	}

	log.Infof("Listing all activities from repository %s", repository)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing activities for pullrequest %s in repository %s with profile %s", listOptions.PullRequestID.Value, repository, profile.Current)) {
		return nil
	}

	activities, err := profile.GetAll[Activity](cmd.Context(), cmd, uripath)
	if err != nil {
		return err
	}
	if len(activities) == 0 {
		log.Infof("No activities found")
		return nil
	}
	core.Sort(activities, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(
		cmd.Context(),
		cmd,
		Activities(core.Filter(activities, func(activity Activity) bool {
			return true
		})),
	)
}
