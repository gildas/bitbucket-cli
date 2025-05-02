package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
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
	Repository string
	Check      bool
}

func init() {
	Command.AddCommand(watchCmd)

	watchCmd.Flags().StringVar(&watchOptions.Repository, "repository", "", "Repository to watch an issue from. Defaults to the current repository")
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
		err = profile.Put(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/watch", args[0]),
			nil,
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to watch issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return
}
