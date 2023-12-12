package profile

import (
	"os"
	"path/filepath"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a profile",
	Args:  cobra.NoArgs,
	RunE:  createProcess,
}

var createOptions struct {
	Profile
	OutputFormat common.EnumFlag
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.OutputFormat = common.EnumFlag{Allowed: []string{"json", "yaml", "table"}, Value: ""}
	createCmd.Flags().StringVarP(&createOptions.Name, "name", "n", "", "Name of the profile")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the profile")
	createCmd.Flags().BoolVar(&createOptions.Default, "default", false, "True if this is the default profile")
	createCmd.Flags().StringVarP(&createOptions.User, "user", "u", "", "User's name of the profile")
	createCmd.Flags().StringVar(&createOptions.Password, "password", "", "Password of the profile")
	createCmd.Flags().StringVar(&createOptions.ClientID, "client-id", "", "Client ID of the profile")
	createCmd.Flags().StringVar(&createOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	createCmd.Flags().StringVar(&createOptions.AccessToken, "access-token", "", "Access Token of the profile")
	createCmd.Flags().Var(&createOptions.OutputFormat, "output", "Output format (json, yaml, table).")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagsRequiredTogether("user", "password")
	createCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	createCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
}

func createProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	log.Infof("Creating profile %s", createOptions.Name)
	if len(createOptions.OutputFormat.String()) > 0 {
		createOptions.Profile.OutputFormat = createOptions.OutputFormat.String()
	}
	if err := createOptions.Validate(); err != nil {
		return err
	}
	if _, found := Profiles.Find(createOptions.Name); found {
		return errors.DuplicateFound.With("name", createOptions.Name)
	}

	Profiles.Add(&createOptions.Profile)
	viper.Set("profiles", Profiles)
	if len(viper.ConfigFileUsed()) > 0 {
		log.Infof("Writing configuration to %s", viper.ConfigFileUsed())
		return viper.WriteConfig()
	}
	if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		configPath := filepath.Join(configDir, "bitbucket")
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return err
		}
		configFile := filepath.Join(configPath, "config-cli.yml")
		if err := viper.WriteConfigAs(configFile); err != nil {
			return err
		}
		if info, err := os.Stat(configFile); err == nil && info.Mode() != 0600 {
			return os.Chmod(configFile, 0600)
		}
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		return viper.WriteConfigAs(filepath.Join(homeDir, ".bitbucket-cli"))
	} else {
		return err
	}
}
