package project

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show", "info", "display"},
	Short:   "get a profile",
	Args:    cobra.ExactArgs(1),
	RunE:    getProcess,
}

var getOptions struct {
	Repository string
	Workspace  string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get pullrequest from. Defaults to the current repository")
	getCmd.Flags().StringVar(&getOptions.Workspace, "workspace", "", "Workspace to get pullrequest from. Defaults to the current workspace")
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying project %s", args[0])
	var project Project

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		getOptions.Repository,
		fmt.Sprintf("/workspaces/%s/projects/%s", getOptions.Workspace, args[0]),
		&project,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project %s: %s\n", args[0], err)
		os.Exit(1)
	}

	payload, _ := json.MarshalIndent(project, "", "  ")
	fmt.Println(string(payload))
	return nil
}
