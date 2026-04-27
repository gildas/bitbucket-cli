package commit

import (
	"fmt"
	"io"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:               "diff [flags] <commit-hash> [<commit-hash>]",
	Short:             "show the diff of a commit or between two commits",
	Args:              cobra.RangeArgs(1, 2),
	ValidArgsFunction: validCommitArgs,
	RunE:              diffProcess,
}

var diffOptions struct {
	Stat bool
}

func init() {
	Command.AddCommand(diffCmd)

	diffCmd.Flags().BoolVar(&diffOptions.Stat, "stat", false, "show only the diffstat")
}

func validCommitArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "diff")
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

func diffProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "diff")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	spec := args[0]
	if len(args) == 2 {
		spec += ".." + args[1]
	}

	log.Debugf("Displaying diff for spec: %s", spec)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing diff for %s", spec)) {
		return nil
	}
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	uripath := repository.GetPath("diff", spec)
	if diffOptions.Stat {
		uripath = repository.GetPath("diffstat", spec)
	}

	diff, err := profile.GetRaw(log.ToContext(cmd.Context()), cmd, uripath)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, diff)
	return err
}
