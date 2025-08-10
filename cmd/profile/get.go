package profile

import (
	"fmt"
	"os"

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
	RunE:              getProcess,
}

var getOptions struct {
	Current bool
	Columns *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)
	getOptions.Columns = flags.NewEnumSliceFlag(columns...)

	getCmd.Flags().BoolVar(&getOptions.Current, "current", false, "Get the current profile")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if getOptions.Current {
		log.Infof("Displaying current profile")
		return Current.Print(cmd.Context(), cmd, Current)
	}

	if len(args) == 0 {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying profile %s (Valid names: %v)", args[0], Profiles.Names())
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
	return Current.Print(cmd.Context(), cmd, profile)
}
