package task

import (
	"fmt"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <task-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a pullrequest task by its <task-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	PullRequestID *flags.EnumFlag
	Columns       *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.PullRequestID = flags.NewEnumFlagWithFunc(getCmd, "", prcommon.GetPullRequestIDs)
	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().Var(getOptions.PullRequestID, "pullrequest", "Pullrequest to get tasks from")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.MarkFlagRequired("pullrequest")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	taskIDs, err := GetPullRequestTaskIDs(cmd.Context(), cmd, getOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return taskIDs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Displaying pullrequest task %s", args[0])
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing pullrequest task %s", args[0])) {
		return nil
	}

	var task Task

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", getOptions.PullRequestID.Value, "tasks", args[0]),
		&task,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pullrequest task %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, task)
}
