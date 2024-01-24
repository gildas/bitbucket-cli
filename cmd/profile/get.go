package profile

import (
	"github.com/gildas/go-errors"
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
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&getOptions.Current, "current", false, "Get the current profile")
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
	return Current.Print(cmd.Context(), cmd, profile)
}
