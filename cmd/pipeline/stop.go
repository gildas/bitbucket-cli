package pipeline

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop [flags] <pipeline-uuid-or-build-number>",
	Aliases: []string{"cancel", "abort"},
	Short:   "stop a running pipeline",
	Args:    cobra.ExactArgs(1),
	RunE:    stopProcess,
}

var stopOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(stopCmd)

	stopCmd.Flags().StringVar(&stopOptions.Repository, "repository", "", "Repository to stop pipeline in. Defaults to the current repository")
}

func stopProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "stop")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	pipelineID := args[0]

	log.Infof("Stopping pipeline %s", pipelineID)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Stopping pipeline %s", pipelineID) {
		return nil
	}

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pipelines/%s/stopPipeline", pipelineID),
		nil,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop pipeline %s: %s\n", pipelineID, err)
		os.Exit(1)
	}

	fmt.Printf("Pipeline %s stopped successfully\n", pipelineID)
	return nil
}
