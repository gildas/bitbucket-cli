package step

import (
	"fmt"
	"strconv"

	plcommon "bitbucket.org/gildas_cherruel/bb/cmd/pipeline/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pipeline steps",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository      string
	PipelineID      *flags.EnumFlag
	Columns         *flags.EnumSliceFlag
	SortBy          *flags.EnumFlag
	ShowLogsCommand bool
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.PipelineID = flags.NewEnumFlagWithFunc("", plcommon.GetPipelineIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pipeline steps from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.PipelineID, "pipeline", "Pipeline to list steps from")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	listCmd.Flags().BoolVar(&listOptions.ShowLogsCommand, "show-logs-command", false, "Show the command to get the logs for each step")
	_ = listCmd.MarkFlagRequired("pipeline")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.PipelineID.CompletionFunc("pipeline"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all comments from repository %s with profile %s", listOptions.Repository, profile.Current)
	steps, err := profile.GetAll[Step](
		cmd.Context(),
		cmd,
		fmt.Sprintf("pipelines/%s/steps", listOptions.PipelineID.Value),
	)
	if err != nil {
		return err
	}
	if len(steps) == 0 {
		log.Infof("No comment found")
		return nil
	}
	core.Sort(steps, columns.SortBy(listOptions.SortBy.Value))
	steps = core.Map(steps, func(step Step) Step {
		step.BuildNumber, _ = strconv.ParseUint(listOptions.PipelineID.Value, 10, 64)
		step.ShowLogsCommand = listOptions.ShowLogsCommand
		log.Debugf("Updated step %s with BuildNumber %d and ShowLogsCommand=%v", step.ID.String(), step.BuildNumber, step.ShowLogsCommand)
		return step
	})
	return profile.Current.Print(cmd.Context(), cmd, Steps(steps))
}
