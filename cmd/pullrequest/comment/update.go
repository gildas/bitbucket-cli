package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type CommentUpdator struct {
	Content ContentUpdator     `json:"content" mapstructure:"content"`
	Anchor  *common.FileAnchor `json:"inline,omitempty" mapstructure:"inline"`
}

type ContentUpdator struct {
	Raw string `json:"raw" mapstructure:"raw"`
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
	PullRequestID *flags.EnumFlag
	Repository    string
	Comment       string
	File          string
	From          int
	To            int
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
	updateCmd.Flags().StringVar(&updateOptions.Repository, "repository", "", "Repository to update a pullrequest comment into. Defaults to the current repository")
	updateCmd.Flags().Var(updateOptions.PullRequestID, "pullrequest", "Pullrequest to update comments to")
	updateCmd.Flags().StringVar(&updateOptions.Comment, "comment", "", "Updated comment of the pullrequest")
	updateCmd.Flags().StringVar(&updateOptions.File, "file", "", "File to comment on")
	updateCmd.Flags().IntVar(&updateOptions.From, "line", 0, "From line to comment on. Cannot be used with --to")
	updateCmd.Flags().IntVar(&updateOptions.From, "from", 0, "From line to comment on. Cannot be used with --line")
	updateCmd.Flags().IntVar(&updateOptions.To, "to", 0, "To line to comment on. Cannot be used with --line")
	updateCmd.MarkFlagsMutuallyExclusive("line", "from")
	updateCmd.MarkFlagsMutuallyExclusive("line", "to")
	_ = updateCmd.MarkFlagRequired("pullrequest")
	_ = updateCmd.MarkFlagRequired("comment")
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.PullRequestID.CompletionFunc("pullrequest"))
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	commentIDs, err := GetPullRequestCommentIDs(cmd.Context(), cmd, deleteOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return commentIDs, cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	payload := CommentUpdator{
		Content: ContentUpdator{Raw: createOptions.Comment},
	}

	if updateOptions.File != "" {
		payload.Anchor = &common.FileAnchor{
			Path: updateOptions.File,
		}
		if updateOptions.From > 0 {
			payload.Anchor.From = uint64(updateOptions.From)
		}
		if updateOptions.To > 0 {
			payload.Anchor.To = uint64(updateOptions.To)
		}
	} else if updateOptions.From > 0 || updateOptions.To > 0 {
		return errors.RuntimeError.With("Cannot specify from/to without a file")
	}

	log.Record("payload", payload).Infof("Updating pullrequest comment")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating comment %s for pullrequest %s", updateOptions.Comment, updateOptions.PullRequestID) {
		return nil
	}
	var comment Comment

	err = profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments/%s", updateOptions.PullRequestID.Value, args[0]),
		payload,
		&comment,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update comment for pullrequest %s: %s\n", updateOptions.PullRequestID.Value, err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, comment)
}
