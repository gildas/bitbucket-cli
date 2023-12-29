package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:               "approve [flags] <pullrequest-id>",
	Short:             "approve a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: approveValidArgs,
	RunE:              approveProcess,
}

var approveOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(approveCmd)

	approveCmd.Flags().StringVar(&approveOptions.Repository, "repository", "", "Repository to approve pullrequest from. Defaults to the current repository")
}

func approveValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return GetPullRequestIDs(cmd.Context(), cmd, approveOptions.Repository, "OPEN"), cobra.ShellCompDirectiveNoFileComp
}

func approveProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "approve")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Approving pullrequest %s", args[0])
	var participant user.Participant

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/approve", args[0]),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to approve pullrequest %s: %s\n", args[0], err)
		return nil
	}
	return profile.Current.Print(cmd.Context(), participant)
}
