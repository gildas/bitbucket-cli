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

var approveCmd = &cobra.Command{
	Use:               "approve",
	Short:             "approve a pullrequest",
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
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Infof("Getting pullrequests for repository %s", approveOptions.Repository)
	var pullrequests struct {
		Values   []PullRequest `json:"values"`
		PageSize int           `json:"pagelen"`
		Size     int           `json:"size"`
		Page     int           `json:"page"`
	}

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		approveOptions.Repository,
		"pullrequests?state=OPEN",
		&pullrequests,
	)
	if err != nil {
		log.Errorf("Failed to get pullrequests for repository %s", approveOptions.Repository, err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var result []string
	for _, pullrequest := range pullrequests.Values {
		result = append(result, fmt.Sprintf("%d", pullrequest.ID))
	}
	return result, cobra.ShellCompDirectiveNoFileComp
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
		approveOptions.Repository,
		fmt.Sprintf("pullrequests/%s/approve", args[0]),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to approve pullrequest %s: %s\n", args[0], err)
		return nil
	}
	data, _ := json.MarshalIndent(participant, "", "  ")
	fmt.Println(string(data))

	return
}
