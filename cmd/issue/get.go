package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get an issue",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
	Changes    bool
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get an issue from. Defaults to the current repository")
	getCmd.Flags().BoolVar(&getOptions.Changes, "changes", false, "Display changes")
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), profile.Current, getOptions.Repository), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if getOptions.Changes {
		log.Infof("Displaying changes for issue %s", args[0])
		changes, err := profile.GetAll[IssueChange](
			log.ToContext(cmd.Context()),
			profile.Current,
			getOptions.Repository,
			fmt.Sprintf("issues/%s/changes", args[0]),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
		return profile.Current.Print(cmd.Context(), IssueChanges(changes))
	}

	log.Infof("Displaying issue %s", args[0])
	var issue Issue

	err = profile.Current.Get(
		log.ToContext(cmd.Context()),
		getOptions.Repository,
		fmt.Sprintf("issues/%s", args[0]),
		&issue,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get issue %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), issue)
}
