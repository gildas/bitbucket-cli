package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"append"},
	Short:   "add a reviewer",
	Args:    cobra.ExactArgs(1),
	PreRunE: disableUnsupportedFlags,
	RunE:    addProcess,
}

var addOptions struct {
	Project *flags.EnumFlag
}

func init() {
	Command.AddCommand(addCmd)

	addOptions.Project = flags.NewEnumFlagWithFunc("", GetProjectKeys)
	addCmd.Flags().Var(addOptions.Project, "project", "Project Key to add reviewers to")
	_ = addCmd.RegisterFlagCompletionFunc(addOptions.Project.CompletionFunc("project"))
	addCmd.SetHelpFunc(hideUnsupportedFlags)
}

func addProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "add")

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

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Adding default reviewer %s to project %s", args[0], project) {
		return nil
	}
	var user user.User

	err = profile.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", workspace, project, args[0]),
		nil,
		&user,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add reviewer: %s\n", err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, user)
}
