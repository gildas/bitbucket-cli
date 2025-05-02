package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type CommentUpdator struct {
	Content common.RenderedText `json:"content" mapstructure:"content"`
}

var updateCmd = &cobra.Command{
	Use:               "update [flags] <comment-id>",
	Aliases:           []string{"edit"},
	Short:             "update an issue comment by its <comment-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
}

var updateOptions struct {
	IssueID    *flags.EnumFlag
	Repository string
	Comment    string
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	updateCmd.Flags().StringVar(&updateOptions.Repository, "repository", "", "Repository to update an issue into. Defaults to the current repository")
	updateCmd.Flags().Var(updateOptions.IssueID, "issue", "Issue to update comments to")
	updateCmd.Flags().StringVar(&updateOptions.Comment, "comment", "", "Updated comment of the issue")
	_ = updateCmd.MarkFlagRequired("issue")
	_ = updateCmd.MarkFlagRequired("comment")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.IssueID.CompletionFunc("issue"))
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	commentIDs, err := GetIssueCommentIDs(cmd.Context(), cmd, profile.Current, deleteOptions.IssueID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(commentIDs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	payload := CommentUpdator{
		Content: common.RenderedText{
			Raw:    updateOptions.Comment,
			Markup: "markdown",
		},
	}

	log.Record("payload", payload).Infof("Updating issue comment")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating comment %s for issue %s", updateOptions.Comment, updateOptions.IssueID) {
		return nil
	}
	var comment Comment

	err = profile.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/comments/%s", updateOptions.IssueID.Value, args[0]),
		payload,
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update comment for issue %s: %s\n", updateOptions.IssueID.Value, err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
