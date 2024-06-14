package profile

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmd = &cobra.Command{
	Use:               "use [flags] <profile-name>",
	Aliases:           []string{"default"},
	Short:             "set the default profile by its <profile-name>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidProfileNames,
	RunE:              useProcess,
}

func init() {
	Command.AddCommand(useCmd)
}

func useProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	log.Infof("Using profile %s (Valid names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}
	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Using profile %s as default", args[0]) {
		Profiles.SetCurrent(profile.Name)
		viper.Set("profiles", Profiles)
		return viper.WriteConfig()
	}
	return nil
}
