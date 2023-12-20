package artifact

import (
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
	Repository string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list artifacts from. Defaults to the current repository")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all projects from repository %s with profile %s", listOptions.Repository, profile.Current)
	artifacts, err := profile.GetAll[Artifact](cmd.Context(), cmd, profile.Current, "downloads")
	if err != nil {
		return err
	}
	if len(artifacts) == 0 {
		log.Infof("No artifact found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), Artifacts(artifacts))
}
