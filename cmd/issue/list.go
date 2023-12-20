package issue

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all issues",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issues from. Defaults to the current repository")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all issues from repository %s with profile %s", listOptions.Repository, profile.Current)
	issues, err := profile.GetAll[Issue](cmd.Context(), cmd, profile.Current, "issues")
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		log.Infof("No issue found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), Issues(issues))
}
