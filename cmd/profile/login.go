package profile

import "github.com/spf13/cobra"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to a profile",
	Args:  cobra.NoArgs,
	RunE:  loginProcess,
}

func init() {
	Command.AddCommand(loginCmd)
}

func loginProcess(cmd *cobra.Command, args []string) error {
	var log = Log.Child(nil, "login")
	var profile = Profiles.Current()

	log.Infof("Logging in to profile %s", profile.Name)
	if err := profile.Validate(); err != nil {
		return err
	}

	if len(profile.AccessToken) > 0 {
		return nil // Access Tokens are always valid
	}

	return nil
}
