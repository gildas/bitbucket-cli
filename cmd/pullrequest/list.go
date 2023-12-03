package pullrequest

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pullrequests",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	State      string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequests from. Defaults to the current repository")
	listCmd.Flags().StringVar(&listOptions.State, "state", "", "Pull request state to fetch. Defaults to \"all\"")
	// TODO: flag state possible values: "all", "open", "closed", "merged"
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if len(listOptions.State) == 0 {
		listOptions.State = "all"
	}

	log.Infof("Listing all pull requests for repository: %s with profile %s", listOptions.Repository, profile.Current)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"commits",
	)
	if err != nil {
		return err
	}
	if len(pullrequests) == 0 {
		log.Infof("No pullrequest found")
		return
	}
	payload, _ := json.MarshalIndent(pullrequests, "", "  ")
	fmt.Println(string(payload))
	return nil
}
