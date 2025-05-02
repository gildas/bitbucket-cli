package artifact

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
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

	log.Infof("Listing all projects from repository %s", listOptions.Repository)
	artifacts, err := profile.GetAll[Artifact](cmd.Context(), cmd, "downloads")
	if err != nil {
		return err
	}
	if len(artifacts) == 0 {
		log.Infof("No artifact found")
		return nil
	}
	core.Sort(artifacts, func(a, b Artifact) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return profile.Current.Print(cmd.Context(), cmd, Artifacts(artifacts))
}
