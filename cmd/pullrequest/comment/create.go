package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	prcommon "bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type CommentCreator struct {
	Content ContentCreator     `json:"content" mapstructure:"content"`
	Anchor  *common.FileAnchor `json:"inline,omitempty" mapstructure:"inline"`
	Parent  *ParentRef         `json:"parent,omitempty" mapstructure:"parent"`
}

type ContentCreator struct {
	Raw string `json:"raw" mapstructure:"raw"`
}

type ParentRef struct {
	ID int64 `json:"id" mapstructure:"id"`
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
	Repository    string
	Comment       string
	File          string
	From          int
	To            int
	ParentID      int64
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create a pullrequest comment into. Defaults to the current repository")
	createCmd.Flags().Var(createOptions.PullRequestID, "pullrequest", "Pullrequest to create comments to")
	createCmd.Flags().StringVar(&createOptions.Comment, "comment", "", "Comment of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.File, "file", "", "File to comment on")
	createCmd.Flags().IntVar(&createOptions.From, "line", 0, "From line to comment on. Cannot be used with --to")
	createCmd.Flags().IntVar(&createOptions.From, "from", 0, "From line to comment on. Cannot be used with --line")
	createCmd.Flags().IntVar(&createOptions.To, "to", 0, "To line to comment on. Cannot be used with --line")
	createCmd.Flags().Int64Var(&createOptions.ParentID, "parent", 0, "Parent comment ID to reply to")
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

	payload := CommentCreator{
		Content: ContentCreator{Raw: createOptions.Comment},
	}

	if createOptions.ParentID > 0 {
		payload.Parent = &ParentRef{ID: createOptions.ParentID}
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

	log.Record("payload", payload).Infof("Creating pullrequest comment")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating comment for pullrequest %s", createOptions.PullRequestID) {
		return nil
	}
	var comment Comment

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments", createOptions.PullRequestID.Value),
		payload,
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create comment for pullrequest %s: %s\n", createOptions.PullRequestID.Value, err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, comment)
}
