package step

import (
	"io"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	plcommon "bitbucket.org/gildas_cherruel/bb/cmd/pipeline/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:               "report [flags] <pipeline-step-uuid-or-name>",
	Aliases:           []string{"report"},
	Short:             "display the report of a pipeline step",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: reportValidArgs,
	RunE:              reportProcess,
}

var reportOptions struct {
	PipelineID *flags.EnumFlag
}

func init() {
	Command.AddCommand(reportCmd)

	reportOptions.PipelineID = flags.NewEnumFlagWithFunc("", plcommon.GetPipelineIDs)
	reportCmd.Flags().Var(reportOptions.PipelineID, "pipeline", "Pipeline to list steps from")
	_ = reportCmd.MarkFlagRequired("pipeline")
	_ = reportCmd.RegisterFlagCompletionFunc(reportOptions.PipelineID.CompletionFunc("pipeline"))
}

func reportValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	stepIDs, err := GetPipelineStepIDs(cmd.Context(), cmd, reportOptions.PipelineID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(stepIDs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func reportProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "getreport")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	report, err := profile.GetRaw(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pipelines", reportOptions.PipelineID.Value, "steps", args[0], "test_reports"),
	)
	if err != nil {
		return errors.Join(errors.Errorf("failed to get test reports for step %s", args[0]), err)
	}

	// Produce the report output
	_, err = io.Copy(os.Stdout, report)
	return err
}
