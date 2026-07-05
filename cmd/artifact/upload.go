package artifact

import (
	"fmt"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
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
	Progress bool
}

func init() {
	Command.AddCommand(uploadCmd)

	uploadCmd.Flags().BoolVar(&uploadOptions.Progress, "progress", false, "Show progress")
}

func uploadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "upload")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(log.ToContext(cmd.Context()), cmd)
	if err != nil {
		return err
	}

	var merr errors.MultiError
	for _, artifactName := range args {
		if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Uploading artifact %s", artifactName) {
			err := profile.Upload(log.ToContext(cmd.Context()), cmd, repository.GetPath("downloads"), artifactName)
			if err != nil {
				if profile.ShouldStopOnError(cmd) {
					return errors.Join(errors.Errorf("Failed to upload artifact %s", artifactName), err)
				} else {
					merr.Append(err)
				}
			}
			log.Infof("Artifact %s uploaded", artifactName)
		}
	}
	if !merr.IsEmpty() && profile.ShouldWarnOnError(cmd) {
		fmt.Fprintf(os.Stderr, "Failed to upload these artifacts: %s\n", merr)
		return nil
	}
	if profile.ShouldIgnoreErrors(cmd) {
		log.Warnf("Failed to upload these artifacts, but ignoring errors: %s", merr)
		return nil
	}
	return merr.AsError()
}
