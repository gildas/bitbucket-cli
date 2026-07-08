package issue

import (
	"fmt"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:               "watch [flags] <issue-id>",
	Short:             "watch an issue by its <issue-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: watchValidArgs,
	RunE:              watchProcess,
}

var watchOptions struct {
	Check bool
}

func init() {
	Command.AddCommand(watchCmd)

	watchCmd.Flags().BoolVar(&watchOptions.Check, "check", false, "Check if the issue is watched")
}

func watchValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func watchProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "watch")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if watchOptions.Check {
		err = profile.Get(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/watch", args[0]),
			nil,
		)
		return
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Watching issue %s", args[0]) {
		err = profile.Put(log.ToContext(cmd.Context()), cmd, repository.GetPath("issues", args[0], "watch"), nil, nil)
		if err != nil {
			return errors.Join(errors.Errorf("Failed to watch issue %s", args[0]), err)
		}
	}
	return
}
