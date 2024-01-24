package artifact

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [flags] <filename...>",
	Aliases:           []string{"remove", "rm"},
	Short:             "delete artifacts by their <filename>.",
	ValidArgsFunction: deleteValidArgs,
	Args:              cobra.MinimumNArgs(1),
	RunE:              deleteProcess,
}

var deleteOptions struct {
	Repository   string
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete artifacts from. Defaults to the current repository")
	deleteCmd.Flags().BoolVar(&deleteOptions.StopOnError, "stop-on-error", false, "Stop on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	deleteCmd.Flags().BoolVar(&deleteOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	deleteCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
}

func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetArtifactNames(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var merr errors.MultiError
	for _, artifactName := range args {
		if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Deleting artifact %s", artifactName) {
			err := profile.Current.Delete(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("downloads/%s", artifactName),
				nil,
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to delete artifact %s: %s\n", artifactName, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Artifact %s deleted", artifactName)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to delete these artifacts: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to delete these artifacts, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
