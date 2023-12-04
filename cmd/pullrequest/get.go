package pullrequest

import (
	"encoding/json"
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
	Short:             "get a profile",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getOpenPullRequests,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get pullrequest from. Defaults to the current repository")
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying pull request %s", args[0])
	var pullrequest PullRequest

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		getOptions.Repository,
		fmt.Sprintf("pullrequests/%s", args[0]),
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}

	payload, _ := json.MarshalIndent(pullrequest, "", "  ")
	fmt.Println(string(payload))
	return nil
}
