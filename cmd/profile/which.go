package profile

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var whichCmd = &cobra.Command{
	Use:     "which",
	Short:   "display the current profile name",
	Args:    cobra.NoArgs,
	PreRunE: disableUnsupportedFlags,
	RunE:    whichProcess,
}

func init() {
	Command.AddCommand(whichCmd)

	whichCmd.SetHelpFunc(hideUnsupportedFlags)
}

func whichProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "which")
	ctx := log.ToContext(cmd.Context())

	profile, err := GetProfileFromCommand(ctx, cmd)
	if errors.Is(err, errors.Empty) || len(Profiles) == 0 {
		if cmd.Flag("stop-on-error").Value.String() == "true" {
			return errors.Errorf("No profiles found")
		}
		common.Verbose(ctx, cmd, "No profile is currently configured")
		return nil
	}
	if err != nil {
		return err
	}

	return profile.Print(ctx, cmd, Current)
}
