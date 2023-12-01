package pullrequest

import (
	"context"
	"encoding/json"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
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
	var log = Log.Child(nil, "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if len(listOptions.State) == 0 {
		listOptions.State = "all"
	}

	log.Infof("Listing all pull requests for repository: %s with profile %s", listOptions.Repository, profile.Current)
	var pullrequests struct {
		Values   []PullRequest `json:"values"`
		PageSize int           `json:"pagelen"`
		Size     int           `json:"size"`
		Page     int           `json:"page"`
	}

	err = profile.Current.Get(
		log.ToContext(context.Background()),
		listOptions.Repository,
		fmt.Sprintf("pullrequests?state=%s", listOptions.State),
		&pullrequests,
	)
	if err != nil {
		return err
	}
	if len(pullrequests.Values) == 0 {
		log.Infof("No pullrequest found")
		return
	}
	payload, _ := json.MarshalIndent(pullrequests, "", "  ")
	fmt.Println(string(payload))
	return nil
}

/*
{"values": [], "pagelen": 10, "size": 0, "page": 1}
*/
