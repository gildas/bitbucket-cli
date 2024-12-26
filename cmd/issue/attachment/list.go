package attachment

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
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
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issue attachments from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.IssueID, "issue", "Issue to list attachments from")
	_ = listCmd.MarkFlagRequired("issue")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.IssueID.CompletionFunc("issue"))
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
	return profile.Current.Print(cmd.Context(), cmd, Attachments(attachments))
}
