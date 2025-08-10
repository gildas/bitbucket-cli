package comment

import (
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all issue comments",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	IssueID    *flags.EnumFlag
	Columns    *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issue comments from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.IssueID, "issue", "Issue to list comments from")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.MarkFlagRequired("issue")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.IssueID.CompletionFunc("issue"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing all comments from repository %s", listOptions.Repository)
	comments, err := profile.GetAll[Comment](
		cmd.Context(),
		cmd,
		fmt.Sprintf("issues/%s/comments", listOptions.IssueID.Value),
	)
	if err != nil {
		return err
	}
	if len(comments) == 0 {
		log.Infof("No comment found")
		return nil
	}
	core.Sort(comments, func(a, b Comment) bool {
		return strings.Compare(strings.ToLower(a.Content.Raw), strings.ToLower(b.Content.Raw)) == -1
	})
	return profile.Current.Print(
		cmd.Context(),
		cmd,
		Comments(core.Filter(comments, func(comment Comment) bool {
			return len(comment.Content.Raw) > 0
		})),
	)
}
