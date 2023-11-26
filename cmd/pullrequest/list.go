package pullrequest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
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

	if len(listOptions.Repository) == 0 {
		return errors.NotImplemented.WithStack()
	}

	if len(listOptions.State) == 0 {
		listOptions.State = "all"
	}

	log.Infof("Listing all pull requests for repository: %s with profile %s", listOptions.Repository, Profile)
	var pullrequests struct {
		Values   []PullRequest `json:"values"`
		PageSize int           `json:"pagelen"`
		Size     int           `json:"size"`
		Page     int           `json:"page"`
	}

	result, err := request.Send(&request.Options{
		Method:        "GET",
		URL:           core.Must(url.Parse(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/pullrequests?state=%s", listOptions.Repository, listOptions.State))),
		Authorization: request.BearerAuthorization(Profile.AccessToken),
		Timeout:       30 * time.Second,
		Logger:        log,
	}, &pullrequests)
	if err != nil {
		return err
	}
	log.Record("result", string(result.Data)).Infof("Result from Bitbucket")

	payload, _ := json.MarshalIndent(pullrequests, "", "  ")
	fmt.Println(string(payload))
	return nil
}

/*
{"values": [], "pagelen": 10, "size": 0, "page": 1}
*/
