package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"append"},
	Short:   "add a reviewer",
	Args:    cobra.ExactArgs(1),
	RunE:    addProcess,
}

var addOptions struct {
	Workspace common.RemoteValueFlag
	Project   common.RemoteValueFlag
}

func init() {
	Command.AddCommand(addCmd)

	addOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	addOptions.Project = common.RemoteValueFlag{AllowedFunc: GetProjectKeys}
	addCmd.Flags().Var(&addOptions.Workspace, "workspace", "Workspace to add reviewers to")
	addCmd.Flags().Var(&addOptions.Project, "project", "Project Key to add reviewers to")
	_ = addCmd.RegisterFlagCompletionFunc("workspace", addOptions.Workspace.CompletionFunc())
	_ = getCmd.RegisterFlagCompletionFunc("project", addOptions.Project.CompletionFunc())
}

func addProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "add")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(addOptions.Workspace.Value) == 0 {
		addOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(addOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}
	if len(addOptions.Project.Value) == 0 {
		addOptions.Project.Value = profile.Current.DefaultProject
		if len(addOptions.Project.Value) == 0 {
			return errors.ArgumentMissing.With("project")
		}
	}

	if !profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Adding default reviewer %s to project %s", args[0], addOptions.Project) {
		return nil
	}
	var user user.User

	err := profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", addOptions.Workspace, addOptions.Project, args[0]),
		nil,
		&user,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add reviewer: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, user)
}
