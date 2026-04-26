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
	Use:               "update [flags] <slug>",
	Short:             "update a repository in a project and a workspace. The project <slug> must be unique in the workspace.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
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

	updateOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceAllowedSlugs)
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

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	slugs, err := GetRepositorySlugs(cmd.Context(), cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(slugs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := GetRepositoryByName(cmd.Context(), cmd, args[0])
	if err != nil {
		return errors.Join(
			errors.Errorf("failed to get repository: %s", args[0]),
			err,
		)
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

	log.Record("payload", payload).Infof("Updating repository %s/%s in project %s", repository.Workspace.Slug, repository.Slug, updateOptions.Project)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating repository %s/%s in project %s", repository.Workspace.Slug, repository.Slug, updateOptions.Project) {
		return nil
	}
	var updated Repository

	err = profile.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/repositories/%s/%s", repository.Workspace.Slug, repository.Slug),
		payload,
		&updated,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update repository %s/%s: %s\n", repository.Workspace.Slug, repository.Slug, err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, updated)
}
