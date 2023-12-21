package artifact

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"add", "create"},
	Short:   "upload an artifact",
	Args:    cobra.ExactArgs(1),
	RunE:    uploadProcess,
}

var uploadOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&uploadOptions.Repository, "repository", "", "Repository to upload artifacts to. Defaults to the current repository")
}

func uploadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "upload")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Uploading artifact %s", args[0])

	err := profile.Current.Upload(
		log.ToContext(cmd.Context()),
		cmd,
		"downloads",
		args[0],
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to upload artifact %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return nil
}
