package commit

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
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

	log.Infof("Listing all branches for repository: %s with profile %s", listOptions.Repository, profile.Current)
	commits, err := profile.GetAll[Commit](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"commits",
	)
	if err != nil {
		return err
	}
	if len(commits) == 0 {
		log.Infof("No branch found")
		return
	}
	return profile.Current.Print(cmd.Context(), Commits(commits))
}
