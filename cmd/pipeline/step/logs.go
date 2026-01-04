package step

import (
	"fmt"
	"io"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	plcommon "bitbucket.org/gildas_cherruel/bb/cmd/pipeline/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:               "logs [flags] <pipeline-step-uuid-or-name>",
	Aliases:           []string{"log"},
	Short:             "display the logs of a pipeline step",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: logValidArgs,
	RunE:              logProcess,
}

var logOptions struct {
	Repository string
	PipelineID *flags.EnumFlag
}

func init() {
	Command.AddCommand(logCmd)

	logOptions.PipelineID = flags.NewEnumFlagWithFunc("", plcommon.GetPipelineIDs)
	logCmd.Flags().StringVar(&logOptions.Repository, "repository", "", "Repository to get pipeline from. Defaults to the current repository")
	logCmd.Flags().Var(logOptions.PipelineID, "pipeline", "Pipeline to list steps from")
	_ = logCmd.MarkFlagRequired("pipeline")
	_ = logCmd.RegisterFlagCompletionFunc(logOptions.PipelineID.CompletionFunc("pipeline"))
}

func logValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	stepIDs, err := GetPipelineStepIDs(cmd.Context(), cmd, logOptions.PipelineID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(stepIDs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func logProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "getlogs")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	steplog, err := profile.Current.GetRaw(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pipelines/%s/steps/%s/log", logOptions.PipelineID.Value, args[0]),
	)
	if err != nil {
		return errors.Join(errors.Errorf("failed to get logs for step %s", args[0]), err)
	}

	// Produce the log output
	_, err = io.Copy(os.Stdout, steplog)
	return err
}
