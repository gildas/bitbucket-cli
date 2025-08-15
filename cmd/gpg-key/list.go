package gpgkey

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all public GPG keys for a given user",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Owner   string
	Columns *flags.EnumSliceFlag
	SortBy  *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Owner, "user", "", "Owner of the keys")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing all GPG keys for %s", listOptions.Owner)
	keys, err := GetGPGKeys(cmd.Context(), cmd)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		log.Infof("No key found")
	}
	core.Sort(keys, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, GPGKeys(keys))
}
