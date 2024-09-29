package user

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show", "info", "display"},
	Short:   "get a user",
	Args:    cobra.ExactArgs(1),
	RunE:    getProcess,
}

var getOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get an issue from. Defaults to the current repository")
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Displaying account %s", args[0])
	var account *Account

	if strings.ToLower(args[0]) == "myself" || strings.ToLower(args[0]) == "me" {
		account, err = GetMe(cmd.Context(), cmd)
	} else {
		account, err = GetAccount(cmd.Context(), cmd, args[0])
	}
	if err != nil {
		return err
	}
	return profile.Print(cmd.Context(), cmd, account)
}
