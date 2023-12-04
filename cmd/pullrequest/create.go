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
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
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

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		createOptions.Repository,
		"pullrequests",
		payload,
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pullrequest: %s\n", err)
		os.Exit(1)
	}
	ref := struct {
		Title string `json:"title"`
		ID    uint64 `json:"id"`
	}{
		Title: pullrequest.Title,
		ID:    pullrequest.ID,
	}
	data, _ := json.MarshalIndent(ref, "", "  ")
	fmt.Println(string(data))

	return
}
