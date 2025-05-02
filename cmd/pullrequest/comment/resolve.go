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

var resolveCmd = &cobra.Command{
	Use:               "resolve [flags] <comment-id>",
	Aliases:           []string{"remove", "rm"},
	Short:             "resolve a pullrequest comment by its <comment-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: resolveValidArgs,
	RunE:              resolveProcess,
}

var resolveOptions struct {
	PullRequestID *flags.EnumFlag
	Repository    string
}

func init() {
	Command.AddCommand(resolveCmd)

	resolveOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
	resolveCmd.Flags().StringVar(&resolveOptions.Repository, "repository", "", "Repository to resolve a pullrequest comment from. Defaults to the current repository")
	resolveCmd.Flags().Var(resolveOptions.PullRequestID, "pullrequest", "Pullrequest to resolve comments from")
	_ = resolveCmd.MarkFlagRequired("pullrequest")
	_ = resolveCmd.RegisterFlagCompletionFunc(resolveOptions.PullRequestID.CompletionFunc("pullrequest"))
}

func resolveValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	commentIDs, err := GetPullRequestCommentIDs(cmd.Context(), cmd, deleteOptions.PullRequestID.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return commentIDs, cobra.ShellCompDirectiveNoFileComp
}

func resolveProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "resolve")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Resolving comment %s from pullrequest %s", args[0], reopenOptions.PullRequestID) {
		return nil
	}

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/comments/%s/resolve", resolveOptions.PullRequestID.Value, args[0]),
		nil,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve pullrequest comment %s: %s\n", args[0], err)
		os.Exit(1)
	}
	log.Infof("Pullrequest comment %s resolved", args[0])
	return nil
}
