package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var unvoteCmd = &cobra.Command{
	Use:               "unvote [flags] <issue-id>",
	Short:             "remove vote for an issue by its <issue-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: unvoteValidArgs,
	RunE:              unvoteProcess,
}

var unvoteOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(unvoteCmd)

	unvoteCmd.Flags().StringVar(&unvoteOptions.Repository, "repository", "", "Repository to unvote an issue from. Defaults to the current repository")
}

func unvoteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ids, err := GetIssueIDs(cmd.Context(), cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func unvoteProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unvote")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unvoting from issue %s", args[0]) {
		err = profile.Delete(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/vote", args[0]),
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unvote issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return
}
