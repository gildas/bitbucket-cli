package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unvoteCmd = &cobra.Command{
	Use:               "unvote",
	Short:             "remove vote for an issue",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: unvoteValidArgs,
	RunE:              unvoteProcess,
}

var unvoteOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(unvoteCmd)

	unvoteCmd.Flags().StringVar(&unvoteOptions.Repository, "repository", "", "Repository to unvote an issue from. Defaults to the current repository")
}

func unvoteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func unvoteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unvote")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("unvote for issue %s", args[0])
	err = profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/vote", args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unvote issue %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return
}
