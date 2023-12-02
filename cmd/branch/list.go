package branch

import (
	"context"
	"encoding/json"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all branches",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list branches from. Defaults to the current repository")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all branches for repository: %s with profile %s", listOptions.Repository, profile.Current)
	var branches struct {
		Values   []Branch `json:"values"`
		PageSize int      `json:"pagelen"`
		Size     int      `json:"size"`
		Page     int      `json:"page"`
	}

	err = profile.Current.Get(
		log.ToContext(context.Background()),
		listOptions.Repository,
		fmt.Sprintf("refs/branches"),
		&branches,
	)
	if err != nil {
		return err
	}
	if len(branches.Values) == 0 {
		log.Infof("No branch found")
		return
	}
	payload, _ := json.MarshalIndent(branches, "", "  ")
	fmt.Println(string(payload))
	return nil
}
