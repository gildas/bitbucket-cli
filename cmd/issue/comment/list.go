package comment

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
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
	Issue      int
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issues from. Defaults to the current repository")
	listCmd.Flags().IntVar(&listOptions.Issue, "issue", 0, "Issue to list comments from. Defaults to the current issue")
	_ = listCmd.MarkFlagRequired("issue")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all comments from repository %s with profile %s", listOptions.Repository, profile.Current)
	comments, err := profile.GetAll[Comment](
		cmd.Context(),
		profile.Current,
		listOptions.Repository,
		fmt.Sprintf("issues/%d/comments", listOptions.Issue),
	)
	if err != nil {
		return err
	}
	if len(comments) == 0 {
		log.Infof("No issue found")
		return nil
	}
	return profile.Current.Print(
		cmd.Context(),
		Comments(core.Filter(comments, func(comment Comment) bool {
			return len(comment.Content.Raw) > 0
		})),
	)
}
