package pullrequest

import (
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pullrequests",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	State      *flags.EnumFlag
	Query      string
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
	PageLength int
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.State = flags.NewEnumFlag("all", "declined", "merged", "+open", "superseded")
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().Var(listOptions.State, "state", "Pull request state to fetch. Defaults to \"open\"")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter pull requests")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	listCmd.Flags().IntVar(&listOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.State.CompletionFunc("state"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var uripath string

	if len(listOptions.Query) > 0 {
		uripath = repository.GetPath(fmt.Sprintf("pullrequests?state=%s&q=%s", url.QueryEscape(strings.ToUpper(listOptions.State.String())), url.QueryEscape(listOptions.Query)))
	} else {
		uripath = repository.GetPath(fmt.Sprintf("pullrequests?state=%s", url.QueryEscape(strings.ToUpper(listOptions.State.String()))))
	}

	log.Infof("Listing %s pull requests for repository: %s", listOptions.State, repository)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing %s pull requests for repository: %s", listOptions.State, repository)) {
		return nil
	}

	pullrequests, err := profile.GetAll[PullRequest](log.ToContext(cmd.Context()), cmd, uripath)
	if err != nil {
		return err
	}
	if len(pullrequests) == 0 {
		log.Infof("No pullrequest found")
		return
	}
	core.Sort(pullrequests, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, PullRequests(pullrequests))
}
