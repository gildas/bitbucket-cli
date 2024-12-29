package repository

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type RepositoryUpdator struct {
	Name        string                    `json:"name,omitempty"                  mapstructure:"name"`
	Description string                    `json:"description,omitempty" mapstructure:"description"`
	Project     *project.ProjectReference `json:"project,omitempty"     mapstructure:"project"`
	IsPrivate   *bool                     `json:"is_private,omitempty"            mapstructure:"is_private"`
	Language    string                    `json:"language,omitempty"    mapstructure:"language"`
	MainBranch  *branch                   `json:"mainbranch,omitempty"  mapstructure:"mainbranch"`
	ForkPolicy  string                    `json:"fork_policy,omitempty" mapstructure:"fork_policy"`
}

var updateCmd = &cobra.Command{
	Use:   "update [flags] <slug>",
	Short: "update a repository in a project and a workspace. The project <slug> must be unique in the workspace.",
	Args:  cobra.ExactArgs(1),
	RunE:  updateProcess,
}

var updateOptions struct {
	Workspace   *flags.EnumFlag
	Project     *flags.EnumFlag
	Name        string
	Description string
	Public      bool
	Private     bool
	Language    string
	MainBranch  string
	ForkPolicy  *flags.EnumFlag
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	updateOptions.Project = flags.NewEnumFlagWithFunc("", project.GetProjectKeys)
	updateOptions.ForkPolicy = flags.NewEnumFlag("allow_forks", "+no_public_forks", "no_forks")
	updateCmd.Flags().Var(updateOptions.Workspace, "workspace", "Workspace to update repositories from")
	updateCmd.Flags().Var(updateOptions.Project, "project", "Project to update repositories from")
	updateCmd.Flags().StringVar(&updateOptions.Name, "name", "", "Name of the repository")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the repository")
	updateCmd.Flags().BoolVar(&updateOptions.Private, "private", false, "make the repository private")
	updateCmd.Flags().BoolVar(&updateOptions.Public, "public", false, "make the repository public")
	updateCmd.Flags().StringVar(&updateOptions.Language, "language", "", "Language of the repository")
	updateCmd.Flags().StringVar(&updateOptions.MainBranch, "main-branch", "", "Main branch of the repository")
	updateCmd.Flags().Var(updateOptions.ForkPolicy, "fork-policy", "Fork policy of the repository. Default: no_public_forks")
	updateCmd.MarkFlagsMutuallyExclusive("private", "public")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.Workspace.CompletionFunc("workspace"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.Project.CompletionFunc("project"))
}

func updateProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(updateOptions.Workspace.Value) == 0 {
		updateOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(updateOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	var private *bool

	if updateOptions.Private {
		private = &updateOptions.Private // => true
	} else if updateOptions.Public {
		private = &updateOptions.Private // => false
	}

	payload := RepositoryUpdator{
		Name:        updateOptions.Name,
		Description: updateOptions.Description,
		IsPrivate:   private,
		Language:    updateOptions.Language,
		ForkPolicy:  updateOptions.ForkPolicy.Value,
	}
	if len(updateOptions.MainBranch) > 0 {
		payload.MainBranch = &branch{Type: "branch", Name: updateOptions.MainBranch}
	}
	if len(updateOptions.Project.Value) > 0 {
		payload.Project = project.NewReference(updateOptions.Project.Value)
	}

	log.Record("payload", payload).Infof("Updating repository %s/%s in project %s", updateOptions.Workspace, updateOptions.Name, updateOptions.Project)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating repository %s/%s in projecct %s", updateOptions.Workspace, updateOptions.Name, updateOptions.Project) {
		return nil
	}
	var repository Repository

	err = profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/repositories/%s/%s", updateOptions.Workspace, args[0]),
		payload,
		&repository,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update repository %s/%s: %s\n", updateOptions.Workspace, args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, repository)
}
