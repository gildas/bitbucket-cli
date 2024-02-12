package artifact

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:               "download [flags] <filename...>",
	Aliases:           []string{"get", "fetch"},
	Short:             "download artifacts by their <filename>.",
	ValidArgsFunction: downloadValidArgs,
	Args:              cobra.MinimumNArgs(1),
	RunE:              getProcess,
}

var downloadOptions struct {
	Repository   string
	Destination  string
	Progress     bool
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&downloadOptions.Repository, "repository", "", "Repository to download artifacts from. Defaults to the current repository")
	downloadCmd.Flags().StringVar(&downloadOptions.Destination, "destination", "", "Destination folder to download the artifact to. Defaults to the current folder")
	downloadCmd.Flags().BoolVar(&downloadOptions.Progress, "progress", false, "Show progress")
	downloadCmd.Flags().BoolVar(&downloadOptions.StopOnError, "stop-on-error", false, "Stop on error")
	downloadCmd.Flags().BoolVar(&downloadOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	downloadCmd.Flags().BoolVar(&downloadOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	downloadCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
	_ = downloadCmd.MarkFlagDirname("destination")
}

func downloadValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetArtifactNames(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "download")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var merr errors.MultiError
	for _, artifactName := range args {
		if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Downloading artifact %s to %s", artifactName, downloadOptions.Destination) {
			err := profile.Current.Download(
				log.ToContext(cmd.Context()),
				cmd,
				fmt.Sprintf("downloads/%s", artifactName),
				downloadOptions.Destination,
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to download artifact %s: %s\n", artifactName, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Artifact %s downloaded", artifactName)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to download these artifacts: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to download these artifacts, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
