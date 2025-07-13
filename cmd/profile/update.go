package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:               "update [flags] <profile-name>",
	Aliases:           []string{"edit"},
	Short:             "update a profile by its <profile-name>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              updateProcess,
}

var updateOptions struct {
	Profile
	DefaultWorkspace *flags.EnumFlag
	DefaultProject   *flags.EnumFlag
	OutputFormat     *flags.EnumFlag
	CloneProtocol    *flags.EnumFlag
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.DefaultWorkspace = flags.NewEnumFlagWithFunc("", getWorkspaceSlugs)
	updateOptions.DefaultProject = flags.NewEnumFlagWithFunc("", getProjectKeys)
	updateOptions.OutputFormat = flags.NewEnumFlag("json", "yaml", "table")
	updateOptions.CloneProtocol = flags.NewEnumFlag("+git", "https", "ssh")
	updateCmd.Flags().StringVarP(&updateOptions.Name, "name", "n", "", "Name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the profile")
	updateCmd.Flags().BoolVar(&updateOptions.Default, "default", false, "True if this is the default profile")
	updateCmd.Flags().StringVar(&updateOptions.VaultKey, "vault-key", "bitbucket-cli", "Vault key to use for storing credentials. Default is bitbucket-cli. On Windows, the Windows Credential Manager will be used, On Linux and macOS, the system keychain will be used.")
	updateCmd.Flags().StringVarP(&updateOptions.User, "user", "u", "", "User's name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Password, "password", "", "Password of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientID, "client-id", "", "Client ID of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	updateCmd.Flags().StringVar(&updateOptions.AccessToken, "access-token", "", "Access Token of the profile")
	updateCmd.Flags().Var(updateOptions.DefaultWorkspace, "default-workspace", "Default workspace of the profile")
	updateCmd.Flags().Var(updateOptions.DefaultProject, "default-project", "Default project of the profile")
	updateCmd.Flags().Var(updateOptions.CloneProtocol, "clone-protocol", "Default protocol to use for cloning repositories. Default is git, can be https, git, or ssh")
	updateCmd.Flags().StringVar(&updateOptions.CloneVaultKey, "clone-vault-key", "bitbucket-cli-clone", "Vault key to use for authentication when cloning with the https protocol. Default is bitbucket-cli-clone. On Windows, the Windows Credential Manager will be used, On Linux and macOS, the system keychain will be used.")
	updateCmd.Flags().StringVar(&updateOptions.CloneVaultUsername, "clone-vault-username", "", "Username to use for authentication when retrieving credentials from the vault.")
	updateCmd.Flags().Var(updateOptions.OutputFormat, "output", "Output format (json, yaml, table).")
	updateCmd.Flags().Var(&updateOptions.ErrorProcessing, "error-processing", "Error processing (StopOnError, WanOnError, IgnoreErrors).")
	updateCmd.Flags().BoolVar(&updateOptions.Progress, "progress", false, "Show progress during upload/download operations.")
	updateCmd.MarkFlagsRequiredTogether("user", "password")
	updateCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	updateCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultWorkspace.CompletionFunc("default-workspace"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultProject.CompletionFunc("default-project"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.CloneProtocol.CompletionFunc("clone-protocol"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.OutputFormat.CompletionFunc("output"))
	_ = updateCmd.RegisterFlagCompletionFunc("error-processing", updateOptions.ErrorProcessing.CompletionFunc())
}

func updateProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if len(updateOptions.DefaultWorkspace.String()) > 0 {
		updateOptions.Profile.DefaultWorkspace = updateOptions.DefaultWorkspace.String()
	}
	if len(updateOptions.DefaultProject.String()) > 0 {
		updateOptions.Profile.DefaultProject = updateOptions.DefaultProject.String()
	}
	if len(updateOptions.OutputFormat.String()) > 0 {
		updateOptions.Profile.OutputFormat = updateOptions.OutputFormat.String()
	}
	if len(updateOptions.CloneProtocol.String()) > 0 {
		updateOptions.Profile.CloneProtocol = updateOptions.CloneProtocol.String()
	}
	log.Infof("Checking if profile %s exists (Valid Names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}

	log.Record("profile", profile).Debugf("Updating profile %s", profile.Name)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating profile %s", profile.Name) {
		return nil
	}

	// We need to check updates to the vault key early, so we can store the client secret and password in the vault if provided
	if !cmd.Flag("vault-key").Changed {
		updateOptions.VaultKey = profile.VaultKey
	}

	if cmd.Flag("client-secret").Changed && len(updateOptions.ClientSecret) > 0 {
		clientID := profile.ClientID
		if cmd.Flag("client-id").Changed && len(updateOptions.ClientID) > 0 {
			clientID = updateOptions.ClientID
		}
		if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, clientID, updateOptions.ClientSecret); err != nil {
			log.Errorf("Failed to store client secret in the vault, the secret will be stored in plain text in the configuration file", err)
			fmt.Fprintf(os.Stderr, "Failed to store client secret in the vault, the secret will be stored in plain text in the configuration file: %s\n", err)
		} else {
			log.Infof("Stored client secret in the vault for %s", clientID)
			updateOptions.ClientSecret = "" // Clear the secret from the profile
		}
	}
	if len(updateOptions.Password) > 0 {
		user := profile.User
		if cmd.Flag("user").Changed && len(updateOptions.User) > 0 {
			user = updateOptions.User
		}
		if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, user, updateOptions.Password); err != nil {
			log.Errorf("Failed to store user password in the vault, the password will be stored in plain text in the configuration file", err)
			fmt.Fprintf(os.Stderr, "Failed to store user password in the vault, the password will be stored in plain text in the configuration file: %s\n", err)
		} else {
			log.Infof("Stored user password in the vault for %s", user)
			updateOptions.Password = "" // Clear the password from the profile
		}
	}
	if cmd.Flag("access-token").Changed && len(updateOptions.AccessToken) > 0 {
		name := profile.Name
		if cmd.Flag("name").Changed && len(updateOptions.Name) > 0 {
			name = updateOptions.Name
		}
		if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, name, updateOptions.AccessToken); err != nil {
			log.Errorf("Failed to store access token in the vault, the token will be stored in plain text in the configuration file", err)
			fmt.Fprintf(os.Stderr, "Failed to store access token in the vault, the token will be stored in plain text in the configuration file: %s\n", err)
		} else {
			log.Infof("Stored access token in the vault for %s", name)
			updateOptions.AccessToken = "" // Clear the access token from the profile
		}
	}

	err := profile.Update(updateOptions.Profile)
	if err != nil {
		return err
	}
	if cmd.Flags().Changed("progress") {
		profile.Progress = updateOptions.Progress
	}
	if updateOptions.Default {
		Profiles.SetCurrent(profile.Name)
	}
	log.Record("profile", profile).Debugf("Updated profile %s", profile.Name)

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
	return Current.Print(cmd.Context(), cmd, profile)
}
