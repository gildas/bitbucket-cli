package profile

import (
	"os"
	"path/filepath"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create a profile",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	Profile
	DefaultWorkspace *flags.EnumFlag
	DefaultProject   *flags.EnumFlag
	OutputFormat     *flags.EnumFlag
	CloneProtocol    *flags.EnumFlag
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.DefaultWorkspace = flags.NewEnumFlagWithFunc("", getWorkspaceSlugs)
	createOptions.DefaultProject = flags.NewEnumFlagWithFunc("", getProjectKeys)
	createOptions.OutputFormat = flags.NewEnumFlag("json", "yaml", "table")
	createOptions.CloneProtocol = flags.NewEnumFlag("+git", "https", "ssh")
	createCmd.Flags().StringVarP(&createOptions.Name, "name", "n", "", "Name of the profile")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the profile")
	createCmd.Flags().BoolVar(&createOptions.Default, "default", false, "True if this is the default profile")
	createCmd.Flags().StringVarP(&createOptions.User, "user", "u", "", "User's name of the profile")
	createCmd.Flags().StringVar(&createOptions.Password, "password", "", "Password of the profile")
	createCmd.Flags().StringVar(&createOptions.ClientID, "client-id", "", "Client ID of the profile")
	createCmd.Flags().StringVar(&createOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	createCmd.Flags().Uint16Var(&createOptions.CallbackPort, "callback-port", 0, "Port to listen to for the Authorization Code Grant")
	createCmd.Flags().StringVar(&createOptions.AccessToken, "access-token", "", "Access Token of the profile")
	createCmd.Flags().Var(createOptions.DefaultWorkspace, "default-workspace", "Default workspace of the profile")
	createCmd.Flags().Var(createOptions.DefaultProject, "default-project", "Default project of the profile")
	createCmd.Flags().Var(createOptions.CloneProtocol, "clone-protocol", "Default protocol to use for cloning repositories. Default is git, can be https, git, or ssh")
	createCmd.Flags().StringVar(&createOptions.CloneVaultKey, "clone-vault-key", "bitbucket-cli", "Vault key to use for authentication when cloning with the https protocol. Default is bitbucket-cli. On Windows, the Windows Credential Manager will be used, On Linux and macOS, the system keychain will be used.")
	createCmd.Flags().StringVar(&createOptions.CloneVaultUsername, "clone-vault-username", "", "Username to use for authentication when retrieving credentials from the vault.")
	createCmd.Flags().Var(createOptions.OutputFormat, "output", "Output format (json, yaml, table).")
	createCmd.Flags().Var(&createOptions.ErrorProcessing, "error-processing", "Error processing (StopOnError, WanOnError, IgnoreErrors).")
	createCmd.Flags().BoolVar(&createOptions.Progress, "progress", false, "Show progress during upload/download operations.")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagsRequiredTogether("user", "password")
	createCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	createCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.CloneProtocol.CompletionFunc("clone-protocol"))
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.OutputFormat.CompletionFunc("output"))
	_ = createCmd.RegisterFlagCompletionFunc("error-processing", createOptions.ErrorProcessing.CompletionFunc())
}

func createProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if len(createOptions.DefaultWorkspace.String()) > 0 {
		createOptions.Profile.DefaultWorkspace = createOptions.DefaultWorkspace.String()
	}
	if len(createOptions.DefaultProject.String()) > 0 {
		createOptions.Profile.DefaultProject = createOptions.DefaultProject.String()
	}
	if len(createOptions.OutputFormat.String()) > 0 {
		createOptions.Profile.OutputFormat = createOptions.OutputFormat.String()
	}
	if len(createOptions.CloneProtocol.String()) > 0 {
		createOptions.Profile.CloneProtocol = createOptions.CloneProtocol.String()
	}
	log.Infof("Creating profile %s", createOptions.Name)
	if err := createOptions.Validate(); err != nil {
		return err
	}
	if _, found := Profiles.Find(createOptions.Name); found {
		return errors.DuplicateFound.With("name", createOptions.Name)
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating profile %s", createOptions.Name) {
		return nil
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
