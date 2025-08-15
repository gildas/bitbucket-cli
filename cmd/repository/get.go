package repository

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <slug_or_uuid>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a repository by its <slug> or <uuid>. With the --forks flag, it will display the forks of the repository.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Workspace *flags.EnumFlag
	ShowForks bool
	Columns   *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().Var(getOptions.Workspace, "workspace", "Workspace to get repositories from")
	getCmd.Flags().BoolVar(&getOptions.ShowForks, "forks", false, "Show the forks of the repository")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Workspace.CompletionFunc("workspace"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	slugs, err := GetRepositorySlugs(cmd.Context(), cmd, getOptions.Workspace.String())
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(slugs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(getOptions.Workspace.Value) == 0 {
		getOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(getOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	if getOptions.ShowForks {
		log.Infof("Displaying forks of repository %s", args[0])
		forks, err := profile.GetAll[Repository](
			cmd.Context(),
			cmd,
			fmt.Sprintf("/repositories/%s/%s/forks", getOptions.Workspace, args[0]),
		)
		if err != nil {
			return err
		}
		if len(forks) == 0 {
			log.Infof("No fork found")
			return nil
		}
		return profile.Current.Print(cmd.Context(), cmd, Repositories(forks))
	}

	log.Infof("Displaying repository %s", args[0])
	repository, err := GetRepository(log.ToContext(cmd.Context()), cmd, profile.Current, getOptions.Workspace.String(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get repository %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, repository)
}
