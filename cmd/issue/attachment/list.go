package attachment

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
	Short: "list all issue attachments",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	IssueID    *flags.EnumFlag
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns...)
	listOptions.SortBy = flags.NewEnumFlag(sortBy...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issue attachments from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.IssueID, "issue", "Issue to list attachments from")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.MarkFlagRequired("issue")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.IssueID.CompletionFunc("issue"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing all attachments from repository %s", listOptions.Repository)
	attachments, err := profile.GetAll[Attachment](
		cmd.Context(),
		cmd,
		fmt.Sprintf("issues/%s/attachments", listOptions.IssueID.Value),
	)
	if err != nil {
		return err
	}
	if len(attachments) == 0 {
		log.Infof("No issue found")
		return nil
	}
	core.Sort(attachments, func(a, b Attachment) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return profile.Current.Print(cmd.Context(), cmd, Attachments(attachments))
}
