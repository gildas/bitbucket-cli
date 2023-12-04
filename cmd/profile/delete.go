package profile

import (
	"errors"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "delete a profile",
	ValidArgsFunction: ValidProfileNames,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	All bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteOptions.All, "all", false, "Delete all profiles")
}

func deleteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")
	var deleted int

	if deleteOptions.All {
		log.Infof("Deleting all profiles")
		deleted = Profiles.Delete(Profiles.Names()...)
	} else if len(args) == 0 {
		return errors.New("accepts 1 arg(s), received 0")
	} else {
		log.Infof("Deleting profiles %s", args)
		deleted = Profiles.Delete(args...)
	}
	log.Infof("Deleted %d profiles", deleted)
	if deleted == 0 {
		return nil
	}
	viper.Set("profiles", Profiles)
	return viper.WriteConfig()
}
