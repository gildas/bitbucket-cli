package profile

import (
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "display the current profile name",
	Args:  cobra.NoArgs,
	RunE:  whichProcess,
}

func init() {
	Command.AddCommand(whichCmd)
}

func whichProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "which")

	return Current.Print(log.ToContext(cmd.Context()), cmd, Current)
}
