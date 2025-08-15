package pullrequest

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
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
	Repository string
	State      *flags.EnumFlag
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.State = flags.NewEnumFlag("all", "declined", "merged", "+open", "superseded")
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequests from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.State, "state", "Pull request state to fetch. Defaults to \"open\"")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.State.CompletionFunc("state"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing %s pull requests for repository: %s", listOptions.State, listOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		cmd,
		"pullrequests?state="+strings.ToUpper(listOptions.State.String()),
	)
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
