package workspace

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"github.com/spf13/cobra"
)

type Workspace struct {
	Type  string     `json:"type"  mapstructure:"type"`
	ID    string     `json:"uuid"  mapstructure:"uuid"`
	Name  string     `json:"name"  mapstructure:"name"`
	Slug  string     `json:"slug"  mapstructure:"slug"`
	Links link.Links `json:"links" mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}
