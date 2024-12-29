package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "cache",
	Short: "Manage the CLI cache",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cache requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}
