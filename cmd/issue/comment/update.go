package comment

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
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
	IssueID *flags.EnumFlag
	Comment string
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.IssueID = flags.NewEnumFlagWithFunc(updateCmd, "", GetIssueIDs)
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

	commentIDs, err := GetIssueCommentIDs(cmd.Context(), cmd, profile.Current, updateOptions.IssueID.Value)
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

	repository, err := repository.GetRepository(cmd.Context(), cmd)
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
		repository.GetPath("issues", updateOptions.IssueID.Value, "comments", args[0]),
		payload,
		&comment,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to update comment for issue %s", updateOptions.IssueID.Value), err)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
