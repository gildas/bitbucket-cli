package project

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
	Short: "list all projects",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Workspace string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Workspace, "workspace", "", "Workspace to list projects from")
	_ = listCmd.MarkFlagRequired("workspace")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all projects from workspace %s with profile %s", listOptions.Workspace, profile.Current)
	projects, err := profile.GetAll[Project](
		cmd.Context(),
		profile.Current,
		"",
		fmt.Sprintf("/workspaces/%s/projects", listOptions.Workspace),
	)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		log.Infof("No project found")
		return nil
	}
	payload, _ := json.MarshalIndent(projects, "", "  ")
	fmt.Println(string(payload))
	return nil
}
