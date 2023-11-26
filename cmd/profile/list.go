package profile

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all profiles",
	Args:  cobra.NoArgs,
	Run:   listProcess,
}

func init() {
	Command.AddCommand(listCmd)
}

func listProcess(cmd *cobra.Command, args []string) {
	var log = Log.Child(nil, "list")

	log.Infof("Listing all profiles")
	if len(Profiles) == 0 {
		log.Infof("No profiles found")
		return
	}
	payload, _ := json.MarshalIndent(Profiles, "", "  ")
	fmt.Println(string(payload))
}
