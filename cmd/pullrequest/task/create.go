package task

import (
	"fmt"
	"strconv"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/pullrequest/comment"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type TaskCreator struct {
	Content   ContentCreator           `json:"content"           mapstructure:"content"`
	Comment   *comment.ParentReference `json:"comment,omitempty" mapstructure:"comment"`
	IsPending bool                     `json:"pending"           mapstructure:"pending"`
}

type ContentCreator struct {
	Raw string `json:"raw" mapstructure:"raw"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create a pullrequest task",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	PullRequestID *flags.EnumFlag
	Content       string
	CommentID     *flags.EnumFlag
	Pending       bool
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.PullRequestID = flags.NewEnumFlagWithFunc(createCmd, "", prcommon.GetPullRequestIDs)
	createOptions.CommentID = flags.NewEnumFlagWithFunc(createCmd, "", comment.GetPullRequestCommentIDs)
	createCmd.Flags().Var(createOptions.PullRequestID, "pullrequest", "Pullrequest to create tasks to")
	createCmd.Flags().StringVar(&createOptions.Content, "content", "", "Content of the task")
	createCmd.Flags().Var(createOptions.CommentID, "comment", "Comment ID to create task on")
	createCmd.Flags().BoolVar(&createOptions.Pending, "pending", false, "Mark the task as pending")
	_ = createCmd.MarkFlagRequired("pullrequest")
	_ = createCmd.MarkFlagRequired("content")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.PullRequestID.CompletionFunc("pullrequest"))
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.CommentID.CompletionFunc("comment"))
}

func createProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	task := TaskCreator{
		Content: ContentCreator{
			Raw: createOptions.Content,
		},
		IsPending: createOptions.Pending,
	}
	if len(createOptions.CommentID.Value) > 0 {
		commentID, err := strconv.ParseInt(createOptions.CommentID.Value, 10, 64)
		if err != nil {
			return errors.Join(errors.Errorf("Failed to parse comment ID %s", createOptions.CommentID.Value), err)
		}
		task.Comment = &comment.ParentReference{
			ID: commentID,
		}
	}

	log.Infof("Creating pullrequest task on pullrequest %s", createOptions.PullRequestID.Value)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Creating pullrequest task on pullrequest %s", createOptions.PullRequestID.Value)) {
		return nil
	}

	var created Task

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", createOptions.PullRequestID.Value, "tasks"),
		task,
		&created,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to create Pull Request Task on Pull Request %s", createOptions.PullRequestID.Value), err)
	}
	return profile.Print(cmd.Context(), cmd, created)
}
