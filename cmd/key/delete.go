package key

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <fingerprints...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete GPG keys by their <fingerprint>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: deleteValidArgs,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Owner string
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteOptions.Owner, "user", "", "Owner of the keys")
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return GetGPGKeyFingerprints(cmd.Context(), cmd), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Deleting GPG keys %v", args)

	for _, fingerprint := range args {
		err := profile.Delete(
			cmd.Context(),
			cmd,
			fmt.Sprintf("/user/gpg-keys/%s", fingerprint),
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
