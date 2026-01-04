package pipeline

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	plcommon "bitbucket.org/gildas_cherruel/bb/cmd/pipeline/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <pipeline-uuid-or-build-number>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a pipeline by its UUID or build number",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
	Columns    *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get pipeline from. Defaults to the current repository")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	ids, err := plcommon.GetPipelineIDs(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying pipeline %s", args[0])
	var pipeline Pipeline

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pipelines/%s", args[0]),
		&pipeline,
	)
	if err != nil {
		return errors.Join(errors.Errorf("failed to get pipeline %s", args[0]), err)
	}

	return profile.Current.Print(cmd.Context(), cmd, pipeline)
}
