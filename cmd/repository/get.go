package repository

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <slug_or_uuid>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a repository by its <slug> or <uuid>. With the --forks flag, it will display the forks of the repository.",
	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: getValidArgs,
	PreRunE:           disableUnsupportedFlags,
	RunE:              getProcess,
}

var getOptions struct {
	ShowForks bool
	Columns   *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().BoolVar(&getOptions.ShowForks, "forks", false, "Show the forks of the repository")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
	getCmd.SetHelpFunc(hideUnsupportedFlags)
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var repository *Repository

	if len(args) == 0 {
		if repository, err = GetRepository(cmd.Context(), cmd); err != nil {
			return errors.Join(
				errors.Errorf("failed to get current repository"),
				err,
			)
		}
	} else {
		if repository, err = GetRepositoryBySlugOrID(cmd.Context(), cmd, args[0]); err != nil {
			return errors.Join(
				errors.Errorf("failed to get repository: %s", args[0]),
				err,
			)
		}
	}

	if getOptions.ShowForks {
		log.Infof("Displaying forks of repository %s", repository.Slug)
		if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing forks of repository %s", repository.Slug)) {
			return nil
		}

		forks, err := repository.GetForks(cmd.Context(), cmd)
		if err != nil {
			return errors.Join(
				errors.Errorf("Failed to get forks of repository %s", repository.Slug),
				err,
			)
		}
		if len(forks) == 0 {
			log.Infof("No fork found")
			return nil
		}
		return profile.Print(cmd.Context(), cmd, Repositories(forks))
	}

	log.Infof("Displaying repository %s", repository.Slug)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing repository %s", repository.Slug)) {
		return nil
	}
	return profile.Print(cmd.Context(), cmd, repository)
}
