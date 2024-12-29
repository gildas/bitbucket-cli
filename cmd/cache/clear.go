package cache

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the cache",
	Args:  cobra.NoArgs,
	Run:   clearProcess,
}

var cache = common.NewCache[any]()

func init() {
	Command.AddCommand(clearCmd)
}

func clearProcess(cmd *cobra.Command, args []string) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "clear")

	log.Infof("Clearing the cache")
	err := cache.Clear()
	if err != nil {
		log.Errorf("Failed to clear the cache: %s", err)
	}
}
