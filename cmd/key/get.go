package key

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a GPG key by its <fingerprint>",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Owner string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Owner, "user", "", "Owner of the key")
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return GetGPGKeyFingerprints(cmd.Context(), cmd), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	owner, err := user.GetUserFromFlags(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Getting GPG key %s", args[0])
	var key *GPGKey

	err = profile.Get(
		cmd.Context(),
		cmd,
		fmt.Sprintf("/users/%s/gpg-keys/%s", owner.ID, args[0]),
		&key,
	)
	if err != nil {
		return err
	}
	return profile.Print(cmd.Context(), cmd, key)
}
