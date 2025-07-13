package profile

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <profile-name>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete a profile by its <profile-name>.",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              deleteProcess,
}

var deleteOptions struct {
	All          bool
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteOptions.All, "all", false, "Delete all profiles")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
}

func deleteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")
	var deleted int

	if deleteOptions.All {
		log.Infof("Deleting all profiles")
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting all profiles") {
			for _, profileName := range Profiles.Names() {
				if profile, found := Profiles.Find(profileName); found {
					log.Infof("Deleting credential for profile %s", profile.Name)
					if len(profile.ClientID) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.ClientID)
						log.Debugf("Deleted client secret for clientID %s from the vault", profile.ClientID)
					} else if len(profile.User) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.User)
						log.Debugf("Deleted user secret for user %s from the vault", profile.User)
					} else if len(profile.Name) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.Name)
						log.Debugf("Deleted name secret for profile %s from the vault", profile.Name)
					}
				}
			}
			deleted = Profiles.Delete(Profiles.Names()...)
		}
	} else {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting profiles %s", strings.Join(args, ", ")) {
			for _, profileName := range args {
				if profile, found := Profiles.Find(profileName); found {
					log.Infof("Deleting credential for profile %s", profile.Name)
					if len(profile.ClientID) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.ClientID)
						log.Debugf("Deleted client secret for clientID %s from the vault", profile.ClientID)
					} else if len(profile.User) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.User)
						log.Debugf("Deleted user secret for user %s from the vault", profile.User)
					} else if len(profile.Name) > 0 {
						_ = profile.DeleteCredentialFromVault(profile.VaultKey, profile.Name)
						log.Debugf("Deleted name secret for profile %s from the vault", profile.Name)
					}
				}
			}
			deleted = Profiles.Delete(args...)
		}
	}
	log.Infof("Deleted %d profiles", deleted)
	if deleted == 0 || cmd.Flag("dry-run").Changed {
		return nil
	}
	viper.Set("profiles", Profiles)
	return viper.WriteConfig()
}
