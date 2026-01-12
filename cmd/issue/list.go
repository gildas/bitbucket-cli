package issue

import (
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all issues",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	States     *flags.EnumSliceFlag
	Query      string
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.States = flags.NewEnumSliceFlagWithAllAllowed("closed", "duplicate", "invalid", "on hold", "+new", "+open", "resolved", "submitted", "wontfix")
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issues from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.States, "state", "State of the issues to list. Can be repeated. One of: all, closed, duplicate, invalid, on hold, new, open, resolved, submitted, wontfix. Default: open, new")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter issues")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.States.CompletionFunc("state"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	filter := ""
	if !core.Contains(listOptions.States.GetSlice(), "all") {
		if states := listOptions.States.GetSlice(); len(states) > 0 {
			filter = "?q="
			for index, state := range states {
				if index > 0 {
					filter += "+OR+"
				}
				filter += fmt.Sprintf(`state="%s"`, strings.ReplaceAll(state, " ", "+"))
			}
		}
	}

	if len(listOptions.Query) > 0 {
		if len(filter) == 0 {
			filter = "?q="
		} else {
			filter += "+AND+"
		}
		filter += url.QueryEscape(listOptions.Query)
	}

	log.Infof("Listing all issues from repository %s with profile %s", listOptions.Repository, profile.Current)
	issues, err := profile.GetAll[Issue](cmd.Context(), cmd, "issues"+filter)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		log.Infof("No issue found")
		return nil
	}
	core.Sort(issues, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, Issues(issues))
}
