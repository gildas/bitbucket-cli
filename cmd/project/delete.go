package project

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)


var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete a project by its key",
	Args:    cobra.ExactArgs(1),
	RunE:    deleteProcess,
}

var deleteOptions struct {
	Workspace string
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteOptions.Workspace, "workspace", "", "Workspace to delete project from")
	_ = deleteCmd.MarkFlagRequired("workspace")
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying project %s", args[0])
	err := profile.Current.Delete(
		log.ToContext(cmd.Context()),
		"",
		fmt.Sprintf("/workspaces/%s/projects/%s", deleteOptions.Workspace, args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete project %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return nil
}
