package artifact

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an artifact by its filename",
	Args:  cobra.ExactArgs(1),
	RunE:  deleteProcess,
}

var deleteOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteOptions.Repository, "repository", "", "Repository to delete artifacts from. Defaults to the current repository")
}

func deleteProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "delete")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Deleting artifact %s from repository %s with profile %s", args[0], listOptions.Repository, profile.Current)
	err := profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("downloads/%s", args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete artifact %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return nil
}
