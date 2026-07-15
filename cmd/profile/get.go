package profile

import (
	"fmt"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <profile-name>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a profile by its <profile-name>.",
	ValidArgsFunction: ValidProfileNames,
	PreRunE:           disableUnsupportedFlags,
	RunE:              getProcess,
}

var getOptions struct {
	Current bool
	Columns *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)
	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)

	getCmd.Flags().BoolVar(&getOptions.Current, "current", false, "Get the current profile")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
	getCmd.SetHelpFunc(hideUnsupportedFlags)
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")
	ctx := log.ToContext(cmd.Context())

	_, err = GetProfileFromCommand(ctx, cmd)
	if errors.Is(err, errors.Empty) || len(Profiles) == 0 {
		if cmd.Flag("stop-on-error").Value.String() == "true" {
			return errors.Errorf("No profiles found")
		}
		common.Verbose(ctx, cmd, "No profile is currently configured")
		return nil
	}
	if err != nil {
		return err
	}

	if getOptions.Current {
		log.Infof("Displaying current profile")
		if Current == nil {
			if cmd.Flag("stop-on-error").Value.String() == "true" {
				return errors.Errorf("There is no profile configured")
			}
			common.Verbose(ctx, cmd, "No profile is currently configured")
			return nil
		}
		return Current.Print(ctx, cmd, Current)
	}

	if len(args) == 0 {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying profile %s (Valid names: %v)", args[0], Profiles.Names())
	if !common.WhatIf(ctx, cmd, fmt.Sprintf("Showing profile %s", args[0])) {
		return nil
	}

	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}
	if err := profile.Validate(); err != nil {
		if cmd.Flag("stop-on-error").Value.String() == "true" {
			return err
		}
		if cmd.Flag("warn-on-error").Value.String() == "true" {
			log.Warnf("Profile %s is not valid: %v", profile.Name, err)
			fmt.Fprintln(os.Stderr, "Profile", profile.Name, "is not valid:", err)
		}
	}
	if len(Profiles) == 1 {
		profile.Default = true
	}
	return profile.Print(ctx, cmd, profile)
}
