package pullrequest

import (
	"fmt"
	"net/url"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// Activities describes a list of Activity
type Activities []Activity

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (activities Activities) GetHeaders(cmd *cobra.Command) []string {
	return Activity{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (activities Activities) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(activities) {
		return []string{}
	}
	return activities[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (activities Activities) Size() int {
	return len(activities)
}

var activitiesCmd = &cobra.Command{
	Use:               "activities",
	Short:             "List all activities of a pullrequest",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: activitiesValidArgs,
	RunE:              activitiesProcess,
}

var activitiesOptions struct {
	Query      string
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
	PageLength int
}

func init() {
	Command.AddCommand(activitiesCmd)

	activitiesOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(activityColumns.Columns()...)
	activitiesOptions.SortBy = flags.NewEnumFlag(activityColumns.Sorters()...)
	activitiesCmd.Flags().StringVar(&activitiesOptions.Query, "query", "", "Query string to filter activities")
	activitiesCmd.Flags().Var(activitiesOptions.Columns, "columns", "Comma-separated list of columns to display")
	activitiesCmd.Flags().Var(activitiesOptions.SortBy, "sort", "Column to sort by")
	activitiesCmd.Flags().IntVar(&activitiesOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
	_ = activitiesCmd.RegisterFlagCompletionFunc(activitiesOptions.Columns.CompletionFunc("columns"))
	_ = activitiesCmd.RegisterFlagCompletionFunc(activitiesOptions.SortBy.CompletionFunc("sort"))
}

func activitiesValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func activitiesProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "activities")

	currentProfile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot merge Pull Request"), err)
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot list activities for Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot list activities for Pull Request"), err)
	}

	uripath := repository.GetPath(fmt.Sprintf("pullrequests/%s/activity", pullRequestID))

	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("%s?q=%s", uripath, url.QueryEscape(listOptions.Query))
	}

	log.Infof("Listing all activities from repository %s with profile %s", repository, currentProfile)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing activities for pullrequest %s in repository %s with profile %s", pullRequestID, repository, currentProfile)) {
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
	core.Sort(activities, activityColumns.SortBy(listOptions.SortBy.Value))
	return currentProfile.Print(
		cmd.Context(),
		cmd,
		Activities(core.Filter(activities, func(activity Activity) bool {
			return true
		})),
	)
}
