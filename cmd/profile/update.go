package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gildas/bitbucket-cli/cmd/common"
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
	PreRunE:           disableUnsupportedFlags,
	RunE:              updateProcess,
}

var updateOptions struct {
	Profile
	DefaultWorkspace *flags.EnumFlag
	DefaultProject   *flags.EnumFlag
	OutputFormat     *flags.EnumFlag
	CloneProtocol    *flags.EnumFlag
	ToVault          bool
	NoVault          bool
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
	if runtime.GOOS != "windows" {
		updateCmd.Flags().StringVar(&updateOptions.VaultKey, "vault-key", "bitbucket-cli", "Vault key to use for storing credentials. Default is bitbucket-cli. On Windows, the Windows Credential Manager will be used, On Linux and macOS, the system keychain will be used.")
	}
	updateCmd.Flags().StringVarP(&updateOptions.User, "user", "u", "", "User's name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Password, "password", "", "Password of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientID, "client-id", "", "Client ID of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	updateCmd.Flags().StringVar(&updateOptions.AccessToken, "access-token", "", "Access Token of the profile")
	updateCmd.Flags().BoolVar(&updateOptions.ToVault, "to-vault", false, "Store credentials in the vault. This will remove any credentials from the profile and store them in the vault. If the vault key is not provided, it will use the existing vault key of the profile or the default vault key if not set.")
	updateCmd.Flags().BoolVar(&updateOptions.NoVault, "no-vault", false, "Do not use a vault for storing credentials")
	updateCmd.Flags().Var(updateOptions.DefaultWorkspace, "default-workspace", "Default workspace of the profile")
	updateCmd.Flags().Var(updateOptions.DefaultProject, "default-project", "Default project of the profile")
	updateCmd.Flags().Var(updateOptions.CloneProtocol, "clone-protocol", "Default protocol to use for cloning repositories. Default is git, can be https, git, or ssh")
	updateCmd.Flags().StringVar(&updateOptions.CloneUser, "clone-user", "", "Username to use when cloning repositories. Default is the username of the profile.")
	updateCmd.Flags().StringVar(&updateOptions.SshKeyFilename, "default-ssh-key-file", "", "Path to the SSH private key file to use when cloning repositories with the ssh protocol.")
	updateCmd.Flags().Var(updateOptions.OutputFormat, "output", "Output format (json, yaml, table).")
	updateCmd.Flags().IntVar(&updateOptions.DefaultPageLength, "default-page-length", 0, "Default number of items per page to retrieve from Bitbucket (Default: 50).")
	updateCmd.Flags().Var(&updateOptions.ErrorProcessing, "error-processing", "Error processing (StopOnError, WanOnError, IgnoreErrors).")
	updateCmd.Flags().BoolVar(&updateOptions.Progress, "progress", false, "Show progress during upload/download operations.")
	updateCmd.MarkFlagsRequiredTogether("user", "password")
	updateCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	updateCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "no-vault")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "access-token")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "client-id")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "client-secret")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "user")
	updateCmd.MarkFlagsMutuallyExclusive("to-vault", "password")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultWorkspace.CompletionFunc("default-workspace"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultProject.CompletionFunc("default-project"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.CloneProtocol.CompletionFunc("clone-protocol"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.OutputFormat.CompletionFunc("output"))
	_ = updateCmd.RegisterFlagCompletionFunc("error-processing", updateOptions.ErrorProcessing.CompletionFunc())
	updateCmd.SetHelpFunc(hideUnsupportedFlags)
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
	log.Infof("Loading profile %s (Valid Names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}

	log.Record("profile", profile).Debugf("Updating profile %s", profile.Name)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating profile %s", profile.Name) {
		return nil
	}

	if updateOptions.ToVault {
		updateOptions.NoVault = false
		vaultKey := profile.VaultKey
		if runtime.GOOS != "windows" && cmd.Flag("vault-key").Changed && len(updateOptions.VaultKey) > 0 {
			vaultKey = updateOptions.VaultKey
		}

		if len(profile.ClientSecret) > 0 {
			if err := profile.SetCredentialInVault(vaultKey, profile.ClientID, profile.ClientSecret); err != nil {
				return errors.Join(errors.Errorf("Failed to store client secret in the vault"), err)
			}
			log.Infof("Stored client secret in the vault for %s", profile.ClientID)
			profile.ClientSecret = ""
			updateOptions.ClientSecret = ""
		} else if len(profile.Password) > 0 {
			if err := profile.SetCredentialInVault(vaultKey, profile.User, profile.Password); err != nil {
				return errors.Join(errors.Errorf("Failed to store user password in the vault"), err)
			}
			log.Infof("Stored user password in the vault for %s", profile.User)
			profile.Password = ""
			updateOptions.Password = ""
		} else if len(profile.AccessToken) > 0 {
			if err := profile.SetCredentialInVault(vaultKey, profile.Name, profile.AccessToken); err != nil {
				return errors.Join(errors.Errorf("Failed to store access token in the vault"), err)
			}
			log.Infof("Stored access token in the vault for %s", profile.Name)
			profile.AccessToken = ""
			updateOptions.AccessToken = ""
		}
	}

	if len(profile.AccessToken) > 0 || len(profile.ClientSecret) > 0 || len(profile.Password) > 0 {
		log.Infof("Profile %s stored its credentials in plain text, we should keep it that way", profile.Name)
		updateOptions.NoVault = true
	}

	// We need to check updates to the vault key early, so we can store the client secret and password in the vault if provided
	if runtime.GOOS != "windows" && !cmd.Flag("vault-key").Changed {
		if len(profile.VaultKey) == 0 {
			profile.VaultKey = "bitbucket-cli"
		}
		updateOptions.VaultKey = profile.VaultKey
	}

	if cmd.Flag("client-secret").Changed && len(updateOptions.ClientSecret) > 0 {
		clientID := profile.ClientID
		if cmd.Flag("client-id").Changed && len(updateOptions.ClientID) > 0 {
			clientID = updateOptions.ClientID
		}
		if !updateOptions.NoVault {
			if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, clientID, updateOptions.ClientSecret); err != nil {
				log.Errorf("Failed to store client secret in the vault, the secret will be stored in plain text in the configuration file", err)
				fmt.Fprintf(os.Stderr, "Failed to store client secret in the vault, the secret will be stored in plain text in the configuration file: %s\n", err)
			} else {
				log.Infof("Stored client secret in the vault for %s", clientID)
				updateOptions.ClientSecret = "" // Clear the secret from the profile
			}
		}
	}
	if cmd.Flag("password").Changed && len(updateOptions.Password) > 0 {
		user := profile.User
		if cmd.Flag("user").Changed && len(updateOptions.User) > 0 {
			user = updateOptions.User
		}
		if !updateOptions.NoVault {
			if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, user, updateOptions.Password); err != nil {
				log.Errorf("Failed to store user password in the vault, the password will be stored in plain text in the configuration file", err)
				fmt.Fprintf(os.Stderr, "Failed to store user password in the vault, the password will be stored in plain text in the configuration file: %s\n", err)
			} else {
				log.Infof("Stored user password in the vault for %s", user)
				updateOptions.Password = "" // Clear the password from the profile
			}
		}
	}
	if cmd.Flag("access-token").Changed && len(updateOptions.AccessToken) > 0 {
		name := profile.Name
		if cmd.Flag("name").Changed && len(updateOptions.Name) > 0 {
			name = updateOptions.Name
		}
		if !updateOptions.NoVault {
			if err := updateOptions.SetCredentialInVault(updateOptions.VaultKey, name, updateOptions.AccessToken); err != nil {
				log.Errorf("Failed to store access token in the vault, the token will be stored in plain text in the configuration file", err)
				fmt.Fprintf(os.Stderr, "Failed to store access token in the vault, the token will be stored in plain text in the configuration file: %s\n", err)
			} else {
				log.Infof("Stored access token in the vault for %s", name)
				updateOptions.AccessToken = "" // Clear the access token from the profile
			}
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
