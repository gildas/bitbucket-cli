package comment

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type CommentCreator struct {
	Content ContentCreator     `json:"content"           mapstructure:"content"`
	Anchor  *common.FileAnchor `json:"inline,omitempty"  mapstructure:"inline"`
	Parent  *ParentReference   `json:"parent,omitempty"  mapstructure:"parent"`
	Pending *bool              `json:"pending,omitempty" mapstructure:"pending"`
}

type ContentCreator struct {
	Raw string `json:"raw" mapstructure:"raw"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create a pullrequest comment",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	PullRequestID *flags.EnumFlag
	Comment       string
	File          string
	From          int
	To            int
	ParentID      int64
	Pending       bool
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.PullRequestID = flags.NewEnumFlagWithFunc(createCmd, "", prcommon.GetPullRequestIDs)
	createCmd.Flags().Var(createOptions.PullRequestID, "pullrequest", "Pullrequest to create comments to")
	createCmd.Flags().StringVar(&createOptions.Comment, "comment", "", "Comment of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.File, "file", "", "File to comment on")
	createCmd.Flags().IntVar(&createOptions.From, "line", 0, "From line to comment on. Cannot be used with --to")
	createCmd.Flags().IntVar(&createOptions.From, "from", 0, "From line to comment on. Cannot be used with --line")
	createCmd.Flags().IntVar(&createOptions.To, "to", 0, "To line to comment on. Cannot be used with --line")
	createCmd.Flags().Int64Var(&createOptions.ParentID, "parent", 0, "Parent comment ID to reply to")
	createCmd.Flags().BoolVar(&createOptions.Pending, "pending", false, "Mark the comment as pending")
	createCmd.MarkFlagsMutuallyExclusive("line", "from")
	createCmd.MarkFlagsMutuallyExclusive("line", "to")
	_ = createCmd.MarkFlagRequired("pullrequest")
	_ = createCmd.MarkFlagRequired("comment")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.PullRequestID.CompletionFunc("pullrequest"))
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
		Content: ContentCreator{Raw: createOptions.Comment},
	}

	if createOptions.ParentID > 0 {
		payload.Parent = &ParentReference{ID: createOptions.ParentID}
	}

	if createOptions.File != "" {
		payload.Anchor = &common.FileAnchor{
			Path: createOptions.File,
		}
		if createOptions.From > 0 {
			payload.Anchor.From = uint64(createOptions.From)
		}
		if createOptions.To > 0 {
			payload.Anchor.To = uint64(createOptions.To)
		}
	} else if createOptions.From > 0 || createOptions.To > 0 {
		return errors.RuntimeError.With("Cannot specify from/to without a file")
	}
	if cmd.Flag("pending").Changed {
		payload.Pending = &createOptions.Pending
	}

	log.Record("payload", payload).Infof("Creating pullrequest comment")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating comment for pullrequest %s", createOptions.PullRequestID.Value) {
		return nil
	}
	var comment Comment

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", createOptions.PullRequestID.Value, "comments"),
		payload,
		&comment,
	)
	if err != nil {
		return errors.Join(errors.Errorf("Failed to create comment for pullrequest %s", createOptions.PullRequestID.Value), err)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
