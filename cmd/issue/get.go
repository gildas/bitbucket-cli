package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <issue-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get an issue by its <issue-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
	Changes    bool
	Columns    *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get an issue from. Defaults to the current repository")
	getCmd.Flags().BoolVar(&getOptions.Changes, "changes", false, "Display changes")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ids, err := GetIssueIDs(cmd.Context(), cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if getOptions.Changes {
		log.Infof("Displaying changes for issue %s", args[0])
		changes, err := GetIssueChanges(cmd.Context(), cmd, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
		return profile.Print(cmd.Context(), cmd, IssueChanges(changes))
	}

	log.Infof("Displaying issue %s", args[0])
	var issue Issue

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s", args[0]),
		&issue,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get issue %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, issue)
}
