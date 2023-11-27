package pullrequest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
	"github.com/spf13/cobra"
)

type PullRequestCreator struct {
	Title             string     `json:"title"`
	Description       string     `json:"description,omitempty"`
	Source            Endpoint   `json:"source"`
	Destination       *Endpoint  `json:"destination,omitempty"`
	Reviewers         []Reviewer `json:"reviewers,omitempty"`
	CloseSourceBranch bool       `json:"close_source_branch,omitempty"`
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a pullrequest",
	Args:  cobra.NoArgs,
	RunE:  createProcess,
}

var createOptions struct {
	Repository        string
	Title             string
	Description       string
	Source            string
	Destination       string
	Reviewers         []string
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create pullrequest from. Defaults to the current repository")
	createCmd.Flags().StringVar(&createOptions.Title, "title", "", "Title of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.Source, "source", "", "Source branch of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.Destination, "destination", "", "Destination branch of the pullrequest")
	createCmd.Flags().StringSliceVar(&createOptions.Reviewers, "reviewer", []string{}, "Reviewer of the pullrequest")
	createCmd.Flags().BoolVar(&createOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch of the pullrequest")
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	var log = Log.Child(nil, "create")

	if len(createOptions.Repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		createOptions.Repository = remote.Repository()
	}

	if len(createOptions.Title) == 0 {
		return errors.ArgumentMissing.With("title")
	}
	if len(createOptions.Source) == 0 {
		return errors.ArgumentMissing.With("source")
	}

	payload := PullRequestCreator{
		Title:             createOptions.Title,
		Description:       createOptions.Description,
		Source:            Endpoint{Branch: Branch{Name: createOptions.Source}},
		CloseSourceBranch: createOptions.CloseSourceBranch,
	}
	if len(createOptions.Destination) > 0 {
		payload.Destination = &Endpoint{Branch: Branch{Name: createOptions.Destination}}
	}

	log.Record("payload", payload).Infof("Creating pullrequest")
	var pullrequest PullRequest

	result, err := request.Send(&request.Options{
		Method:        "POST",
		URL:           core.Must(url.Parse(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/pullrequests", listOptions.Repository))),
		Authorization: request.BearerAuthorization(Profile.AccessToken),
		Timeout:       30 * time.Second,
		Logger:        log,
	}, &pullrequest)
	if err != nil {
		return err
	}
	log.Record("result", string(result.Data)).Infof("Result from Bitbucket")
	data, _ := json.MarshalIndent(pullrequest, "", "  ")
	fmt.Println(string(data))

	return nil
}
