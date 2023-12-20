package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:               "vote",
	Short:             "vote for an issue",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: voteValidArgs,
	RunE:              voteProcess,
}

var voteOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(voteCmd)

	voteCmd.Flags().StringVar(&voteOptions.Repository, "repository", "", "Repository to vote an issue from. Defaults to the current repository")
}

func voteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func voteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "vote")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Vote for issue %s", args[0])
	err = profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/vote", args[0]),
		nil,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to vote issue %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return
}
