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

type RepositoryCreator struct {
	Name        string                    `json:"name"                  mapstructure:"name"`
	Description string                    `json:"description,omitempty" mapstructure:"description"`
	Project     *project.ProjectReference `json:"project,omitempty"     mapstructure:"project"`
	IsPrivate   bool                      `json:"is_private"            mapstructure:"is_private"`
	Language    string                    `json:"language,omitempty"    mapstructure:"language"`
	MainBranch  *branch                   `json:"mainbranch,omitempty"  mapstructure:"mainbranch"`
	ForkPolicy  string                    `json:"fork_policy,omitempty" mapstructure:"fork_policy"`
}

var createCmd = &cobra.Command{
	Use:   "create [flags] <slug>",
	Short: "create a repository in a project and a workspace. The repository <slug> must be unique in the workspace.",
	Args:  cobra.ExactArgs(1),
	RunE:  createProcess,
}

var createOptions struct {
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
	Command.AddCommand(createCmd)

	createOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	createOptions.Project = flags.NewEnumFlagWithFunc("", project.GetProjectKeys)
	createOptions.ForkPolicy = flags.NewEnumFlag("allow_forks", "+no_public_forks", "no_forks")
	createCmd.Flags().Var(createOptions.Workspace, "workspace", "Workspace to create repositories from")
	createCmd.Flags().Var(createOptions.Project, "project", "Project to create repositories from")
	createCmd.Flags().StringVar(&createOptions.Name, "name", "", "Name of the repository")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the repository")
	createCmd.Flags().BoolVar(&createOptions.Private, "private", false, "make the repository private")
	createCmd.Flags().BoolVar(&createOptions.Public, "public", false, "make the repository public")
	createCmd.Flags().StringVar(&createOptions.Language, "language", "", "Language of the repository")
	createCmd.Flags().StringVar(&createOptions.MainBranch, "main-branch", "", "Main branch of the repository")
	createCmd.Flags().Var(createOptions.ForkPolicy, "fork-policy", "Fork policy of the repository. Default: no_public_forks")
	createCmd.MarkFlagsMutuallyExclusive("private", "public")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Workspace.CompletionFunc("workspace"))
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Project.CompletionFunc("project"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(createOptions.Workspace.Value) == 0 {
		createOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(createOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	payload := RepositoryCreator{
		Name:        createOptions.Name,
		Description: createOptions.Description,
		IsPrivate:   createOptions.Private,
		Language:    createOptions.Language,
		ForkPolicy:  createOptions.ForkPolicy.Value,
	}
	if len(createOptions.MainBranch) > 0 {
		payload.MainBranch = &branch{Type: "branch", Name: createOptions.MainBranch}
	}
	if len(createOptions.Project.Value) > 0 {
		payload.Project = project.NewReference(createOptions.Project.Value)
	}

	log.Record("payload", payload).Infof("Creating repository %s/%s in project %s", createOptions.Workspace, createOptions.Name, createOptions.Project)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating repository %s/%s in project %s", createOptions.Workspace, createOptions.Name, createOptions.Project) {
		return nil
	}
	var repository Repository

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/repositories/%s/%s", createOptions.Workspace, args[0]),
		payload,
		&repository,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create repository %s/%s: %s\n", createOptions.Workspace, args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, repository)
}
