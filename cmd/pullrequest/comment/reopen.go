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

var reopenCmd = &cobra.Command{
	Use:               "reopen [flags] <comment-id>",
	Aliases:           []string{"remove", "rm"},
	Short:             "reopen a pullrequest comment by its <comment-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: reopenValidArgs,
	RunE:              reopenProcess,
}

var reopenOptions struct {
	PullRequestID *flags.EnumFlag
	Repository    string
}

func init() {
	Command.AddCommand(reopenCmd)

	reopenOptions.PullRequestID = flags.NewEnumFlagWithFunc("", GetPullRequestIDs)
	reopenCmd.Flags().StringVar(&reopenOptions.Repository, "repository", "", "Repository to reopen a pullrequest comment from. Defaults to the current repository")
	reopenCmd.Flags().Var(reopenOptions.PullRequestID, "pullrequest", "Pullrequest to reopen comments from")
	_ = reopenCmd.MarkFlagRequired("pullrequest")
	_ = reopenCmd.RegisterFlagCompletionFunc("pullrequest", reopenOptions.PullRequestID.CompletionFunc("pullrequest"))
}

func reopenValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetPullRequestCommentIDs(cmd.Context(), cmd, profile.Current, reopenOptions.PullRequestID.Value), cobra.ShellCompDirectiveNoFileComp
}

func reopenProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "reopen")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Resolving comment %s from pullrequest %s", args[0], reopenOptions.PullRequestID) {
		return nil
	}

	err = profile.Current.Delete(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments/%s/resolve", reopenOptions.PullRequestID.Value, args[0]),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reopen pullrequest comment %s: %s\n", args[0], err)
		os.Exit(1)
	}
	log.Infof("Pullrequest comment %s reopened", args[0])
	return nil
}
