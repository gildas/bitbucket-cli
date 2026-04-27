package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
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

func init() {
	Command.AddCommand(unwatchCmd)
}

func unwatchValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func unwatchProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "unwatch")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Unwatching issue %s", args[0]) {
		err = profile.Delete(log.ToContext(cmd.Context()), cmd, repository.GetPath("issues", args[0], "watch"), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unwatch issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return
}
