package comment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
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
}

func init() {
	Command.AddCommand(resolveCmd)

	resolveOptions.PullRequestID = flags.NewEnumFlagWithFunc("", prcommon.GetPullRequestIDs)
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

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Resolving comment %s from pullrequest %s", args[0], resolveOptions.PullRequestID) {
		return nil
	}

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		repository.GetPath("pullrequests", resolveOptions.PullRequestID.Value, "comments", args[0], "resolve"),
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
