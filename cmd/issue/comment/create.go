package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
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
	IssueID *flags.EnumFlag
	Comment string
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	createCmd.Flags().Var(createOptions.IssueID, "issue", "Issue to create comments to")
	createCmd.Flags().StringVar(&createOptions.Comment, "comment", "", "Comment of the issue")
	_ = createCmd.MarkFlagRequired("issue")
	_ = createCmd.MarkFlagRequired("comment")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.IssueID.CompletionFunc("issue"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
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

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("issues", createOptions.IssueID.Value, "comments"),
		payload,
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create comment for issue %s: %s\n", createOptions.IssueID.Value, err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
