package step

import (
	"io"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	plcommon "github.com/gildas/bitbucket-cli/cmd/pipeline/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var casesCmd = &cobra.Command{
	Use:               "cases [flags] <pipeline-step-uuid-or-name>",
	Aliases:           []string{"cases"},
	Short:             "list the test cases of a pipeline step",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: casesValidArgs,
	RunE:              casesProcess,
}

var casesOptions struct {
	PipelineID *flags.EnumFlag
}

func init() {
	Command.AddCommand(casesCmd)

	casesOptions.PipelineID = flags.NewEnumFlagWithFunc(casesCmd, "", plcommon.GetPipelineIDs)
	casesCmd.Flags().Var(casesOptions.PipelineID, "pipeline", "Pipeline to list steps from")
	_ = casesCmd.MarkFlagRequired("pipeline")
	_ = casesCmd.RegisterFlagCompletionFunc(casesOptions.PipelineID.CompletionFunc("pipeline"))
}

func casesValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	stepIDs, err := GetPipelineStepIDs(cmd.Context(), cmd, casesOptions.PipelineID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(stepIDs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func casesProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "listcases")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	cases, err := profile.GetRaw(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pipelines", casesOptions.PipelineID.Value, "steps", args[0], "test_reports", "test_cases"),
	)
	if err != nil {
		return errors.Join(errors.Errorf("failed to get test cases for step %s", args[0]), err)
	}

	// Produce the cases output
	_, err = io.Copy(os.Stdout, cases)
	return err
}
