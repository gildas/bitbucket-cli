package commit

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var ancestorCmd = &cobra.Command{
	Use:               "ancestor <commit-hash> <commit-hash>",
	Short:             "show the ancestor commit of two commits",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: validAncestorArgs,
	RunE:              ancestorProcess,
}

func init() {
	Command.AddCommand(ancestorCmd)
}

func validAncestorArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "ancestor")
	if len(args) == 0 && toComplete == "" || len(args) > 2 {
		log.Debugf("No args or too many args for completion: %v, toComplete: %s", args, toComplete)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Debugf("Getting commit hashes for completion with args: %v and toComplete: %s", args, toComplete)
	names, err := GetCommitHashes(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(names, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func ancestorProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "ancestor")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Debugf("Displaying ancestor for commit %s and %s", args[0], args[1])
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing ancestor for commit %s and %s", args[0], args[1])) {
		return nil
	}
	var ancestor Commit

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("merge-base/%s..%s", args[0], args[1]),
		&ancestor,
	)
	if err != nil {
		return err
	}

	return profile.Print(cmd.Context(), cmd, ancestor)
}
