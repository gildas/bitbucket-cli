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
	ValidArgsFunction: getOpenPullRequests,
	RunE:              declineProcess,
}

var declineOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(declineCmd)

	declineCmd.Flags().StringVar(&declineOptions.Repository, "repository", "", "Repository to decline pullrequest from. Defaults to the current repository")
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
		os.Exit(1)
	}
	data, _ := json.MarshalIndent(participant, "", "  ")
	fmt.Println(string(data))

	return
}
