package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unapproveCmd = &cobra.Command{
	Use:               "unapprove",
	Short:             "unapprove a pullrequest",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: unapproveValidArgs,
	RunE:              unapproveProcess,
}

var unapproveOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(unapproveCmd)

	unapproveCmd.Flags().StringVar(&unapproveOptions.Repository, "repository", "", "Repository to unapprove pullrequest from. Defaults to the current repository")
}

func unapproveValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Infof("Getting open pullrequests for repository %s", approveOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"pullrequests?state=OPEN",
	)
	if err != nil {
		log.Errorf("Failed to get pullrequests for repository %s", unapproveOptions.Repository, err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var result []string
	for _, pullrequest := range pullrequests {
		result = append(result, fmt.Sprintf("%d", pullrequest.ID))
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

func unapproveProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unapprove")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Unapproving pullrequest %s", args[0])
	err = profile.Current.Delete(
		log.ToContext(cmd.Context()),
		unapproveOptions.Repository,
		fmt.Sprintf("pullrequests/%s/approve", args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unapprove pullrequest %s: %s\n", args[0], err)
		return nil
	}
	return
}
