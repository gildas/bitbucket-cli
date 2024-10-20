package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <comment-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get an issue comment by its <comment-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	IssueID    *flags.EnumFlag
	Repository string
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get an issue comment from. Defaults to the current repository")
	getCmd.Flags().Var(getOptions.IssueID, "issue", "Issue to get comments from")
	_ = getCmd.MarkFlagRequired("issue")
	_ = getCmd.RegisterFlagCompletionFunc("issue", getOptions.IssueID.CompletionFunc("issue"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	commentIDs, err := GetIssueCommentIDs(cmd.Context(), cmd, profile.Current, deleteOptions.IssueID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return commentIDs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Displaying issue comment %s", args[0])
	var comment Comment

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/comments/%s", getOptions.IssueID.Value, args[0]),
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get issue comment %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
