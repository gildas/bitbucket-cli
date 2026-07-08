package reviewer

import (
	"fmt"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/user"
	"github.com/gildas/bitbucket-cli/cmd/workspace"
	errors "github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a reviewer",
	ValidArgsFunction: getValidArgs,
	Args:              cobra.ExactArgs(1),
	PreRunE:           disableUnsupportedFlags,
	RunE:              getProcess,
}

var getOptions struct {
	Project *flags.EnumFlag
	Columns *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Project = flags.NewEnumFlagWithFunc(getCmd, "", GetProjectKeys)
	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().Var(getOptions.Project, "project", "Project Key to get reviewers from")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Project.CompletionFunc("project"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
	getCmd.SetHelpFunc(hideUnsupportedFlags)
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	userIDs, err := GetReviewerUserIDs(cmd.Context(), cmd, getOptions.Project.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return userIDs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, err := workspace.GetWorkspace(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	project, err := GetProjectName(cmd, profile)
	if err != nil {
		return err
	}

	log.Infof("Displaying reviewer %s", args[0])
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing reviewer %s", args[0])) {
		return nil
	}
	var user user.User

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", workspace, project, args[0]),
		&user,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to get reviewer %s", args[0]), err)
	}
	return profile.Print(cmd.Context(), cmd, user)
}
