package commit

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
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
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if len(listOptions.Repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		listOptions.Repository = remote.Repository()
	}

	log.Infof("Listing all branches for repository: %s with profile %s", listOptions.Repository, profile.Current)
	var commits struct {
		Values   []Commit `json:"values"`
		PageSize int      `json:"pagelen"`
		Size     int      `json:"size"`
		Page     int      `json:"page"`
	}

	err = profile.Current.Get(
		log.ToContext(context.Background()),
		listOptions.Repository,
		"commits",
		&commits,
	)
	if err != nil {
		return err
	}
	if len(commits.Values) == 0 {
		log.Infof("No branch found")
		return
	}
	payload, _ := json.MarshalIndent(commits, "", "  ")
	fmt.Println(string(payload))
	return nil
}
