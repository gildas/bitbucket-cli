package profile

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a profile",
	Args:  cobra.NoArgs,
	RunE:  createProcess,
}

var createOptions Profile

func init() {
	Command.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createOptions.Name, "name", "n", "", "Name of the profile")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the profile")
	createCmd.Flags().BoolVar(&createOptions.Default, "default", false, "True if this is the default profile")
	createCmd.Flags().StringVarP(&createOptions.User, "user", "u", "", "User's name of the profile")
	createCmd.Flags().StringVar(&createOptions.Password, "password", "", "Password of the profile")
	createCmd.Flags().StringVar(&createOptions.AccessToken, "access-token", "", "Access Token of the profile")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagsRequiredTogether("user", "password")
	createCmd.MarkFlagsMutuallyExclusive("user", "access-token")
}

func createProcess(cmd *cobra.Command, args []string) error {
	log := Log.Child(nil, "create")

	log.Infof("Creating profile %s", createOptions.Name)
	if err := createOptions.Validate(); err != nil {
		return err
	}

	Profiles.Add(&createOptions)
	viper.Set("profiles", Profiles)
	return viper.WriteConfig()
}
