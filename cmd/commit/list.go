package commit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-request"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all commits",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list commits from. Defaults to the current repository")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	var log = Log.Child(nil, "list")

	if len(listOptions.Repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		listOptions.Repository = remote.Repository()
	}

	log.Infof("Listing all branches for repository: %s with profile %s", listOptions.Repository, Profile)
	var commits struct {
		Values   []Commit `json:"values"`
		PageSize int      `json:"pagelen"`
		Size     int      `json:"size"`
		Page     int      `json:"page"`
	}

	result, err := request.Send(&request.Options{
		Method:        "GET",
		URL:           core.Must(url.Parse(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/commits", listOptions.Repository))),
		Authorization: request.BearerAuthorization(Profile.AccessToken),
		Timeout:       30 * time.Second,
		Logger:        log,
	}, &commits)
	if err != nil {
		return err
	}
	log.Record("result", string(result.Data)).Infof("Result from Bitbucket")
	if len(commits.Values) == 0 {
		log.Infof("No branch found")
		return
	}
	payload, _ := json.MarshalIndent(commits, "", "  ")
	fmt.Println(string(payload))
	return nil
}
