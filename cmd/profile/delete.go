package profile

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <profile-name>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete a profile by its <profile-name>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	All          bool
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteOptions.All, "all", false, "Delete all profiles")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
}

func deleteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")
	var deleted int

	if deleteOptions.All {
		log.Infof("Deleting all profiles")
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting all profiles") {
			deleted = Profiles.Delete(Profiles.Names()...)
		}
	} else {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting profiles %s", strings.Join(args, ", ")) {
			deleted = Profiles.Delete(args...)
		}
	}
	log.Infof("Deleted %d profiles", deleted)
	if deleted == 0 || cmd.Flag("dry-run").Changed {
		return nil
	}
	viper.Set("profiles", Profiles)
	return viper.WriteConfig()
}
