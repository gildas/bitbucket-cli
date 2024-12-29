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
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.DefaultWorkspace = flags.NewEnumFlagWithFunc("", getWorkspaceSlugs)
	updateOptions.DefaultProject = flags.NewEnumFlagWithFunc("", getProjectKeys)
	updateOptions.OutputFormat = flags.NewEnumFlag("json", "yaml", "table")
	updateCmd.Flags().StringVarP(&updateOptions.Name, "name", "n", "", "Name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the profile")
	updateCmd.Flags().BoolVar(&updateOptions.Default, "default", false, "True if this is the default profile")
	updateCmd.Flags().StringVarP(&updateOptions.User, "user", "u", "", "User's name of the profile")
	updateCmd.Flags().StringVar(&updateOptions.Password, "password", "", "Password of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientID, "client-id", "", "Client ID of the profile")
	updateCmd.Flags().StringVar(&updateOptions.ClientSecret, "client-secret", "", "Client Secret of the profile")
	updateCmd.Flags().StringVar(&updateOptions.AccessToken, "access-token", "", "Access Token of the profile")
	updateCmd.Flags().Var(updateOptions.DefaultWorkspace, "default-workspace", "Default workspace of the profile")
	updateCmd.Flags().Var(updateOptions.DefaultProject, "default-project", "Default project of the profile")
	updateCmd.Flags().Var(updateOptions.OutputFormat, "output", "Output format (json, yaml, table).")
	updateCmd.Flags().Var(&updateOptions.ErrorProcessing, "error-processing", "Error processing (StopOnError, WanOnError, IgnoreErrors).")
	updateCmd.Flags().BoolVar(&updateOptions.Progress, "progress", false, "Show progress during upload/download operations.")
	updateCmd.MarkFlagsRequiredTogether("user", "password")
	updateCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	updateCmd.MarkFlagsMutuallyExclusive("user", "client-id", "access-token")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultWorkspace.CompletionFunc("default-workspace"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.DefaultProject.CompletionFunc("default-project"))
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
	log.Infof("Checking if profile %s exists (Valid Names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}

	log.Record("profile", profile).Debugf("Updating profile %s", profile.Name)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating profile %s", profile.Name) {
		return nil
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
