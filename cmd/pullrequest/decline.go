package pullrequest

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var declineCmd = &cobra.Command{
	Use:               "decline",
	Short:             "decline a pullrequest",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: declineValidArgs,
	RunE:              declineProcess,
}

var declineOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(declineCmd)

	declineCmd.Flags().StringVar(&declineOptions.Repository, "repository", "", "Repository to decline pullrequest from. Defaults to the current repository")
}

func declineValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		log.Errorf("Failed to get pullrequests for repository %s", declineOptions.Repository, err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var result []string
	for _, pullrequest := range pullrequests {
		result = append(result, fmt.Sprintf("%d", pullrequest.ID))
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

func declineProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "decline")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Declining pullrequest %s", args[0])
	var participant user.Participant

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		declineOptions.Repository,
		fmt.Sprintf("pullrequests/%s/decline", args[0]),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decline pullrequest %s: %s\n", args[0], err)
		return nil
	}
	data, _ := json.MarshalIndent(participant, "", "  ")
	fmt.Println(string(data))

	return
}
