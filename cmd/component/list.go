package component

import (
	"fmt"
	"net/url"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all components",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	Query      string
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list components from. Defaults to the current repository")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter components")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	uripath := "components"
	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("components?q=%s", url.QueryEscape(listOptions.Query))
	}

	log.Infof("Listing all issues from repository %s", listOptions.Repository)
	components, err := profile.GetAll[Component](cmd.Context(), cmd, uripath)
	if err != nil {
		return err
	}
	if len(components) == 0 {
		log.Infof("No component found")
		return nil
	}
	core.Sort(components, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, Components(components))
}
