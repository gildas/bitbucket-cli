package issue

import (
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
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

func init() {
	Command.AddCommand(unvoteCmd)
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

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unvoting from issue %s", args[0]) {
		err = profile.Delete(log.ToContext(cmd.Context()), cmd, repository.GetPath("issues", args[0], "vote"), nil)
		if err != nil {
			return errors.Join(errors.Errorf("Failed to unvote issue %s", args[0]), err)
		}
	}
	return
}
