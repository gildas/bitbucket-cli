package profile

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all profiles",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

func init() {
	Command.AddCommand(listCmd)
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(Command.Name(), "list")

	log.Infof("Listing all profiles")
	if len(Profiles) == 0 {
		log.Infof("No profiles found")
		return
	}
	core.Sort(Profiles, func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return Current.Print(cmd.Context(), cmd, Profiles)
}
