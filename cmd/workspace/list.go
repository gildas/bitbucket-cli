package workspace

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
	Short: "list all workspaces for the current user",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	WithMembership bool
	Query          string
	Columns        *flags.EnumSliceFlag
	SortBy         *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().BoolVar(&listOptions.WithMembership, "membership", false, "List also the workspace memberships of the current user")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter workspaces")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if listOptions.WithMembership {
		log.Infof("Listing all workspace memberships for current user")
		memberships, err := profile.GetAll[Membership](cmd.Context(), cmd, "/user/permissions/workspaces")
		if err != nil {
			return err
		}
		if len(memberships) == 0 {
			log.Infof("No workspace found")
			return nil
		}
		return profile.Current.Print(cmd.Context(), cmd, Memberships(memberships))
	}

	uripath := "workspaces"
	if len(listOptions.Query) > 0 {
		uripath = uripath + "?q=" + listOptions.Query
		uripath = fmt.Sprintf("workspaces?q=%s", url.QueryEscape(listOptions.Query))
	}

	log.Infof("Listing all workspaces")
	workspaces, err := profile.GetAll[Workspace](cmd.Context(), cmd, uripath)
	if err != nil {
		return err
	}
	if len(workspaces) == 0 {
		log.Infof("No workspace found")
		return nil
	}
	core.Sort(workspaces, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, Workspaces(workspaces))
}
