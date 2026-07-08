package profile

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all profiles",
	Args:    cobra.NoArgs,
	PreRunE: disableUnsupportedFlags,
	RunE:    listProcess,
}

var listOptions struct {
	Columns *flags.EnumSliceFlag
	SortBy  *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
	listCmd.SetHelpFunc(hideUnsupportedFlags)
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(Command.Name(), "list")
	ctx := log.ToContext(cmd.Context())

	log.Infof("Listing all profiles")
	if !common.WhatIf(ctx, cmd, "Showing profiles") {
		return nil
	}

	profile, err := GetProfileFromCommand(ctx, cmd)
	if errors.Is(err, errors.Empty) || len(Profiles) == 0 {
		if cmd.Flag("stop-on-error").Value.String() == "true" {
			return errors.Errorf("No profiles found")
		}
		common.Verbose(ctx, cmd, "No profiles found")
		return nil
	}
	if err != nil {
		return err
	}
	core.Sort(Profiles, columns.SortBy(listOptions.SortBy.Value))
	Profiles = core.Map(Profiles, func(profile *Profile) *Profile {
		_ = profile.Validate()
		return profile
	})
	return profile.Print(ctx, cmd, Profiles)
}
