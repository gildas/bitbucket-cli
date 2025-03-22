package sshkey

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all public SSH keys for a given user",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Owner string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Owner, "user", "", "Owner of the keys")
}

func listProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing all SSH keys for %s", listOptions.Owner)
	keys, err := GetSSHKeys(cmd.Context(), cmd)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		log.Infof("No key found")
	}
	return profile.Current.Print(cmd.Context(), cmd, SSHKeys(keys))
}
