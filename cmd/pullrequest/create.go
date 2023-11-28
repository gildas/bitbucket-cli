package pullrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gildas/go-errors"
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

	err = Profile.Post(
		log.ToContext(context.Background()),
		createOptions.Repository,
		"pullrequests",
		payload,
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pullrequest: %s\n", err)
		return nil
	}
	data, _ := json.MarshalIndent(pullrequest, "", "  ")
	fmt.Println(string(data))

	return
}
