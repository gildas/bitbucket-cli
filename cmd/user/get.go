package user

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
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
	Columns    *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns...)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get an issue from. Defaults to the current repository")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Displaying user %s", args[0])
	user, err := GetUser(cmd.Context(), cmd, args[0])
	if err != nil {
		return err
	}
	log.Record("user", user).Debugf("User %s retrieved", args[0])
	return profile.Print(cmd.Context(), cmd, user)
}
