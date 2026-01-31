package permission

import (
	"fmt"
	"net/url"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	wkcommon "bitbucket.org/gildas_cherruel/bb/cmd/workspace/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:               "list",
	Short:             "list all workspace permissions",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: listValidArgs,
	RunE:              listProcess,
}

var listOptions struct {
	Query      string
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
	PageLength int
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter comments")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	listCmd.Flags().IntVar(&listOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
	_ = listCmd.MarkFlagRequired("workspace")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	slugs, err := wkcommon.GetWorkspaceSlugs(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return slugs, cobra.ShellCompDirectiveNoFileComp
}

func listProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var uripath string

	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("/workspaces/%s/permissions?q=%s", args[0], url.QueryEscape(listOptions.Query))
	} else {
		uripath = fmt.Sprintf("/workspaces/%s/permissions", args[0])
	}

	log.Infof("Listing all permissions from workspace %s with profile %s", args[0], profile.Current)
	permissions, err := profile.GetAll[Permission](cmd.Context(), cmd, uripath)
	if err != nil {
		return err
	}
	if len(permissions) == 0 {
		log.Infof("No permission found")
		return nil
	}
	core.Sort(permissions, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, Permissions(permissions))
}
