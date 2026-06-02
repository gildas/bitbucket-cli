package task

import (
	"fmt"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type TaskUpdator struct {
	Content *ContentUpdator `json:"content,omitempty" mapstructure:"content,omitempty"`
	State   string          `json:"state,omitempty"   mapstructure:"state,omitempty"`
}

type ContentUpdator struct {
	Raw string `json:"raw" mapstructure:"raw"`
}

var updateCmd = &cobra.Command{
	Use:               "update [flags] <task-id>",
	Aliases:           []string{"edit"},
	Short:             "update a pullrequest task by its <task-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
}

var updateOptions struct {
	PullRequestID *flags.EnumFlag
	Content       string
	State         *flags.EnumFlag
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.PullRequestID = flags.NewEnumFlagWithFunc(updateCmd, "", prcommon.GetPullRequestIDs)
	updateOptions.State = flags.NewEnumFlag("RESOLVED", "UNRESOLVED")
	updateCmd.Flags().Var(updateOptions.PullRequestID, "pullrequest", "Pullrequest to update tasks to")
	updateCmd.Flags().StringVar(&updateOptions.Content, "content", "", "Updated content of the task")
	updateCmd.Flags().Var(updateOptions.State, "state", "Updated state of the task. Can be one of RESOLVED or UNRESOLVED")
	_ = updateCmd.MarkFlagRequired("pullrequest")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.State.CompletionFunc("state"))
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	taskIDs, err := GetPullRequestTaskIDs(cmd.Context(), cmd, updateOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return taskIDs, cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	taskID := args[0]

	taskUpdator := TaskUpdator{}
	if len(updateOptions.Content) > 0 {
		taskUpdator.Content = &ContentUpdator{
			Raw: updateOptions.Content,
		}
	}
	if updateOptions.State != nil && len(updateOptions.State.Value) > 0 {
		taskUpdator.State = updateOptions.State.Value
	}

	log.Infof("Updating pullrequest task %s on pullrequest %s", taskID, updateOptions.PullRequestID.Value)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Updating pullrequest task %s on pullrequest %s", taskID, updateOptions.PullRequestID.Value)) {
		return nil
	}

	var updated Task

	err = profile.Put(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", updateOptions.PullRequestID.Value, "tasks", taskID),
		taskUpdator,
		&updated,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to update Pull Request Task %s on Pull Request %s", taskID, updateOptions.PullRequestID.Value), err)
	}
	return profile.Print(cmd.Context(), cmd, updated)
}
