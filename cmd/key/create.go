package key

import (
	"fmt"
	"io"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "Add a new GPG key",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	User    string `json:"-"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	KeyFile string `json:"-"`
}

func init() {
	Command.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createOptions.User, "user", "", "Owner's User ID of the key, defaults to the current user")
	createCmd.Flags().StringVar(&createOptions.Name, "name", "", "Name of the key")
	createCmd.Flags().StringVar(&createOptions.Key, "key", "", "GPG key to add")
	createCmd.Flags().StringVar(&createOptions.KeyFile, "key-file", "", "File containing the GPG key to add. Use '-' to read from stdin")
	_ = createCmd.MarkFlagFilename("key-file")
	createCmd.MarkFlagsMutuallyExclusive("key", "key-file")
	createCmd.MarkFlagsOneRequired("key", "key-file")
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if len(createOptions.KeyFile) > 0 {
		var data []byte
		if createOptions.KeyFile == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(createOptions.KeyFile)
		}
		if err != nil {
			return err
		}
		createOptions.Key = string(data)
	}

	if len(createOptions.Key) == 0 {
		return errors.ArgumentMissing.With("key")
	}

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	owner, err := user.GetUserFromFlags(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating GPG key for %s", owner) {
		return nil
	}
	log.Infof("Creating GPG key for %s", owner)
	var key *GPGKey

	err = profile.Post(
		cmd.Context(),
		cmd,
		fmt.Sprintf("/users/%s/gpg-keys", owner.ID.String()),
		createOptions,
		&key,
	)
	if err != nil {
		return err
	}

	return profile.Print(cmd.Context(), cmd, key)
}
