package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unwatchCmd = &cobra.Command{
	Use:               "unwatch [flags] <issue-id>",
	Short:             "stop watching an issue by its <issue-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: unwatchValidArgs,
	RunE:              unwatchProcess,
}

var unwatchOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(unwatchCmd)

	unwatchCmd.Flags().StringVar(&unwatchOptions.Repository, "repository", "", "Repository to unwatch an issue from. Defaults to the current repository")
}

func unwatchValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func unwatchProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unwatch")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unwatching issue %s", args[0]) {
		err = profile.Current.Delete(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/watch", args[0]),
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unwatch issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return
}
