package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type CommentCreator struct {
	Content common.RenderedText `json:"content" mapstructure:"content"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create an issue comment",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	IssueID    *flags.EnumFlag
	Repository string
	Comment    string
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create an issue comment into. Defaults to the current repository")
	createCmd.Flags().Var(createOptions.IssueID, "issue", "Issue to create comments to")
	createCmd.Flags().StringVar(&createOptions.Comment, "comment", "", "Comment of the issue")
	_ = createCmd.MarkFlagRequired("issue")
	_ = createCmd.MarkFlagRequired("comment")
	_ = createCmd.RegisterFlagCompletionFunc("issue", createOptions.IssueID.CompletionFunc("issue"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	payload := CommentCreator{
		Content: common.RenderedText{
			Raw:    createOptions.Comment,
			Markup: "markdown",
		},
	}

	log.Record("payload", payload).Infof("Creating issue comment")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating comment for issue %s", createOptions.IssueID) {
		return nil
	}
	var comment Comment

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/comments", createOptions.IssueID.Value),
		payload,
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create comment for issue %s: %s\n", createOptions.IssueID.Value, err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, comment)
}
