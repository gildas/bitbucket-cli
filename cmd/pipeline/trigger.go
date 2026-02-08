package pipeline

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/branch"
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// TriggerBody represents the body for triggering a pipeline
type TriggerBody struct {
	Target    Target     `json:"target"`
	Variables []Variable `json:"variables,omitempty"`
}

var triggerCmd = &cobra.Command{
	Use:     "trigger",
	Aliases: []string{"run", "start", "create"},
	Short:   "trigger a new pipeline",
	Args:    cobra.NoArgs,
	RunE:    triggerProcess,
}

var triggerOptions struct {
	Repository string
	Branch     string
	Tag        string
	Commit     string
	Pattern    string
	Variables  []string
}

func init() {
	Command.AddCommand(triggerCmd)

	triggerCmd.Flags().StringVar(&triggerOptions.Repository, "repository", "", "Repository to trigger pipeline in. Defaults to the current repository")
	triggerCmd.Flags().StringVar(&triggerOptions.Branch, "branch", "", "Branch to run the pipeline on")
	triggerCmd.Flags().StringVar(&triggerOptions.Tag, "tag", "", "Tag to run the pipeline on")
	triggerCmd.Flags().StringVar(&triggerOptions.Commit, "commit", "", "Specific commit hash to run the pipeline on")
	triggerCmd.Flags().StringVar(&triggerOptions.Pattern, "pattern", "", "Custom pipeline pattern to run (e.g., 'deploy-to-prod')")
	triggerCmd.Flags().StringArrayVar(&triggerOptions.Variables, "variable", []string{}, "Pipeline variable in KEY=VALUE format. Can be specified multiple times")
}

func triggerProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "trigger")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	// Build the target
	target := ReferenceTarget{
		Type: "pipeline_ref_target",
	}

	if len(triggerOptions.Tag) > 0 {
		target.ReferenceType = "tag"
		target.ReferenceName = triggerOptions.Tag
	} else if len(triggerOptions.Branch) > 0 {
		target.ReferenceType = "branch"
		target.ReferenceName = triggerOptions.Branch
	} else {
		// Try to detect the current git branch
		currentBranch, err := branch.GetCurrentBranch()
		if err != nil {
			log.Errorf("Failed to retrieve the current branch", err)
			return errors.ArgumentMissing.With("branch or tag", "use --branch or --tag to specify the target")
		}
		target.ReferenceType = currentBranch.GetType()
		target.ReferenceName = currentBranch.Name
		log.Infof("Using current branch: %s", currentBranch.Name)
	}

	if len(triggerOptions.Commit) > 0 {
		target.Commit = commit.CommitReference{Hash: triggerOptions.Commit}
	}

	if len(triggerOptions.Pattern) > 0 {
		target.Selector = &common.Selector{
			Type:    "custom",
			Pattern: triggerOptions.Pattern,
		}
	}

	payload := TriggerBody{
		Target: target,
	}

	// Parse variables
	if len(triggerOptions.Variables) > 0 {
		payload.Variables = make([]Variable, 0, len(triggerOptions.Variables))
		for _, variable := range triggerOptions.Variables {
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) != 2 {
				return errors.ArgumentInvalid.With("variable", variable)
			}
			payload.Variables = append(payload.Variables, Variable{
				Key:   parts[0],
				Value: parts[1],
			})
		}
	}

	log.Record("payload", payload).Infof("Triggering pipeline")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Triggering pipeline") {
		return nil
	}

	var pipeline Pipeline

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		"pipelines/",
		payload,
		&pipeline,
	)
	if err != nil {
		return errors.Join(errors.Errorf("failed to trigger pipeline"), err)
	}

	return profile.Print(cmd.Context(), cmd, pipeline)
}
