package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a profile",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              getProcess,
}

func init() {
	Command.AddCommand(getCmd)
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	log.Infof("Displaying profile %s (Valid names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}
	payload, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(payload))
	return nil
}
