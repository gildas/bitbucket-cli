package artifact

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:     "upload [flags] <filename...>",
	Aliases: []string{"add", "create"},
	Short:   "upload artifacts",
	Args:    cobra.MinimumNArgs(1),
	RunE:    uploadProcess,
}

var uploadOptions struct {
	Repository   string
	Progress     bool
	StopOnError  bool
	WarnOnError  bool
	IgnoreErrors bool
}

func init() {
	Command.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&uploadOptions.Repository, "repository", "", "Repository to upload artifacts to. Defaults to the current repository")
	uploadCmd.Flags().BoolVar(&uploadOptions.StopOnError, "stop-on-error", false, "Stop on error")
	uploadCmd.Flags().BoolVar(&uploadOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	uploadCmd.Flags().BoolVar(&uploadOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	uploadCmd.Flags().BoolVar(&uploadOptions.Progress, "progress", false, "Show progress")
	uploadCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
}

func uploadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "upload")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var merr errors.MultiError
	for _, artifactName := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Uploading artifact %s to %s", artifactName, downloadOptions.Destination) {
			err := profile.Current.Upload(
				log.ToContext(cmd.Context()),
				cmd,
				"downloads",
				args[0],
			)
			if err != nil {
				if profile.Current.ShouldStopOnError(cmd) {
					fmt.Fprintf(os.Stderr, "Failed to upload artifact %s: %s\n", artifactName, err)
					os.Exit(1)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Artifact %s downloaded", artifactName)
		}
	}
	if !merr.IsEmpty() && profile.Current.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to upload these artifacts: %s\n", merr)
		return nil
	}
	if profile.Current.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to upload these artifacts, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
