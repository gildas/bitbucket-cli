package profile

import (
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmd = &cobra.Command{
	Use:       "use",
	Short:     "set the default profile",
	Args:      cobra.ExactArgs(1),
	ValidArgs: Profiles.Names(),
	RunE:      useProcess,
}

func init() {
	Command.AddCommand(useCmd)
}

func useProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	log.Infof("Using profile %s", args[0])
	log.Warnf("Valid names: %s", Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}
	Profiles.SetCurrent(profile.Name)
	viper.Set("profiles", Profiles)
	return viper.WriteConfig()
}
