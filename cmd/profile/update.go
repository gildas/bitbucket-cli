package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "update a profile",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              updateProcess,
}

var updateOptions Profile

func init() {
	Command.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateOptions.Name, "name", "n", "", "Name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the profile")
	updateCmd.Flags().BoolVar(&updateOptions.Default, "default", false, "True if this is the default profile")
	updateCmd.Flags().StringVarP(&updateOptions.User, "user", "u", "", "User's name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Password, "password", "", "Password of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientID, "client-id", "", "Client ID of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	updateCmd.Flags().StringVar(&updateOptions.AccessToken, "access-token", "", "Access Token of the profile")
	updateCmd.MarkFlagsRequiredTogether("user", "password")
	updateCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	updateCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
}

func updateProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	log.Infof("Updating profile %s", createOptions.Name)
	if err := createOptions.Validate(); err != nil {
		return err
	}

	log.Warnf("Valid names: %s", Profiles.Names())

	if _, found := Profiles.Find(args[0]); !found {
		return errors.NotFound.With("profile", args[0])
	}

	Profiles.Add(&updateOptions)
	_ = Profiles.Delete(args[0])

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
		if err := viper.WriteConfigAs(filepath.Join(homeDir, ".bitbucket-cli")); err != nil {
			return err
		}
	} else {
		return err
	}
	profile, _ := Profiles.Find(args[0])
	payload, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(payload))
	return nil
}
